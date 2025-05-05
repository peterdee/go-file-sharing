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
	"file-sharing/utilities"
)

func main() {
	envError := godotenv.Load()
	if envError != nil {
		log.Fatal(envError)
	}

	if _, uploadsError := os.Stat("uploads"); os.IsNotExist(uploadsError) {
		os.Mkdir("uploads", 0755)
	}

	cache.Connect()
	database.Connect()

	http.HandleFunc("GET /", handlers.IndexHandler)
	http.HandleFunc("GET /api", handlers.IndexHandler)
	http.HandleFunc("GET /api/download/{id}", handlers.DownloadHandler)
	http.HandleFunc("GET /api/info/{id}", handlers.InfoHandler)
	http.HandleFunc("POST /api/upload", handlers.UploadHandler)

	port := utilities.GetEnv(constants.ENV_NAMES.Port, constants.DEFAULT_PORT)
	listener, listenError := net.Listen("tcp", ":"+port)
	if listenError != nil {
		log.Fatal(listenError)
	}

	log.Printf("Server is running on port %s", port)

	if serveError := http.Serve(listener, nil); serveError != nil {
		log.Fatal(serveError)
	}
}
