package main

import (
	"log"
	"net"
	"net/http"

	"github.com/joho/godotenv"

	"file-exchange/constants"
	"file-exchange/handlers"
	"file-exchange/utilities"
)

func main() {
	envError := godotenv.Load()
	if envError != nil {
		log.Fatal(envError)
	}

	http.HandleFunc("GET /api/x", handlers.IndexHandler)
	http.HandleFunc("GET /api/", handlers.IndexHandler)

	port := utilities.GetEnv(constants.ENV_NAMES.Port, constants.DEFAULT_PORT)
	listener, listenError := net.Listen("tcp", ":"+port)
	if listenError != nil {
		log.Fatal(listenError)
	}

	log.Printf("Running the server on port %s", port)

	if serveError := http.Serve(listener, nil); serveError != nil {
		log.Fatal(serveError)
	}
}
