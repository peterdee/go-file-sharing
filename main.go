package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/handlers"
	"file-sharing/handlers/manage"
	"file-sharing/middlewares"
	scheduledtasks "file-sharing/scheduled-tasks"
	"file-sharing/utilities"
)

func main() {
	envError := godotenv.Load()
	if envError != nil {
		log.Fatal(envError)
	}

	uploadsDirectoryName := utilities.GetEnv(
		constants.ENV_NAMES.UplaodsDirectoryName,
		constants.DEFAULT_UPLOADS_DIRECTORY_NAME,
	)
	if _, uploadsError := os.Stat(uploadsDirectoryName); os.IsNotExist(uploadsError) {
		fsError := os.Mkdir(uploadsDirectoryName, 0755)
		if fsError != nil {
			log.Fatal(fsError)
		}
	}

	cache.Connect()
	database.Connect()

	scheduledtasks.MarkAsDeleted()

	// public mux
	publicHandlers := http.NewServeMux()
	publicHandlers.HandleFunc("GET /", handlers.IndexHandler)
	publicHandlers.HandleFunc("GET /api", handlers.IndexHandler)
	publicHandlers.HandleFunc("GET /api/download/{id}", handlers.DownloadHandler)
	publicHandlers.HandleFunc("GET /api/info/{id}", handlers.InfoHandler)
	publicHandlers.HandleFunc("POST /api/upload", handlers.UploadHandler)

	// managing mux
	managingHandlers := http.NewServeMux()
	managingHandlers.HandleFunc("GET /files", manage.ListFilesHandler)

	port := utilities.GetEnv(constants.ENV_NAMES.Port, constants.DEFAULT_PORT)
	listener, listenError := net.Listen("tcp", ":"+port)
	if listenError != nil {
		log.Fatal(listenError)
	}

	log.Printf("Server is running on port %s", port)

	rootMux := http.NewServeMux()
	rootMux.Handle("/", publicHandlers)
	rootMux.Handle("/api/manage", managingHandlers)
	serveError := http.Serve(listener, middlewares.WithLogger(rootMux))
	if serveError != nil {
		log.Fatal(serveError)
	}
}
