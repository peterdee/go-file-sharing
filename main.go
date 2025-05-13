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
	"file-sharing/handlers/auth"
	"file-sharing/handlers/manage"
	"file-sharing/handlers/public"
	"file-sharing/handlers/root"
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
	publicHandlers.HandleFunc("GET /", public.IndexHandler)
	publicHandlers.HandleFunc("GET /api", public.IndexHandler)
	publicHandlers.HandleFunc("GET /api/download/{id}", public.DownloadHandler)
	publicHandlers.HandleFunc("GET /api/info/{id}", public.InfoHandler)
	publicHandlers.HandleFunc("POST /api/upload", public.UploadHandler)

	// auth mux
	authHandlers := http.NewServeMux()
	authHandlers.HandleFunc("POST /set-up", auth.SetUpHandler)
	authHandlers.HandleFunc("POST /sign-in", auth.SignInHandler)

	// managing mux
	managingHandlers := http.NewServeMux()
	managingHandlers.HandleFunc("DELETE /file/{id}", manage.DeleteFileHandler)
	managingHandlers.HandleFunc("GET /file/{id}", manage.DeleteFileHandler)
	managingHandlers.HandleFunc("GET /list", manage.ListFilesHandler)

	// root mux
	rootHandlers := http.NewServeMux()
	rootHandlers.HandleFunc("DELETE /account/{id}", root.DeleteAccountHandler)
	rootHandlers.HandleFunc("GET /account/{id}", root.GetAccountHandler)
	rootHandlers.HandleFunc("GET /list", root.ListAccountsHandler)
	rootHandlers.HandleFunc("PATCH /{id}", root.UpdateAccountHandler)
	rootHandlers.HandleFunc("POST /", root.CreateAccountHandler)

	port := utilities.GetEnv(constants.ENV_NAMES.Port, constants.DEFAULT_PORT)
	listener, listenError := net.Listen("tcp", ":"+port)
	if listenError != nil {
		log.Fatal(listenError)
	}

	log.Printf("Server is running on port %s", port)

	rootMux := http.NewServeMux()
	rootMux.Handle("/", publicHandlers)
	rootMux.Handle("/api/auth", authHandlers)
	rootMux.Handle("/api/manage", managingHandlers)
	rootMux.Handle("/api/root", rootHandlers)
	serveError := http.Serve(listener, middlewares.WithLogger(rootMux))
	if serveError != nil {
		log.Fatal(serveError)
	}
}
