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
	"file-sharing/handlers/user"
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

	// auth mux
	authHandlers := http.NewServeMux()
	authHandlers.HandleFunc("POST /set-up", auth.SetUpHandler)
	authHandlers.HandleFunc("POST /sign-in", auth.SignInHandler)

	// managing mux
	managingHandlers := http.NewServeMux()
	managingHandlers.HandleFunc("DELETE /file/{id}", manage.DeleteFileHandler)
	managingHandlers.HandleFunc("GET /file/{id}", manage.DeleteFileHandler)
	managingHandlers.HandleFunc("GET /list", manage.ListFilesHandler)

	// public mux
	publicHandlers := http.NewServeMux()
	publicHandlers.HandleFunc("GET /", public.IndexHandler)
	publicHandlers.HandleFunc("GET /download/{id}", public.DownloadHandler)
	publicHandlers.HandleFunc("GET /info/{id}", public.InfoHandler)
	publicHandlers.HandleFunc("POST /upload", public.UploadHandler)

	// root mux
	rootHandlers := http.NewServeMux()
	rootHandlers.HandleFunc("DELETE /user/{id}", root.DeleteUserHandler)
	rootHandlers.HandleFunc("GET /user/{id}", root.GetUserHandler)
	rootHandlers.HandleFunc("GET /users", root.ListUsersHandler)
	rootHandlers.HandleFunc("PATCH /user/{id}", root.UpdateUserHandler)
	rootHandlers.HandleFunc("POST /user", root.CreateUserHandler)

	// user mux
	userHandlers := http.NewServeMux()
	userHandlers.HandleFunc("GET /", user.GetUserHandler)
	userHandlers.HandleFunc("PATCH /password", user.ChangePasswordHandler)

	port := utilities.GetEnv(constants.ENV_NAMES.Port, constants.DEFAULT_PORT)
	listener, listenError := net.Listen("tcp", ":"+port)
	if listenError != nil {
		log.Fatal(listenError)
	}

	log.Printf("Server is running on port %s", port)

	combineMux := http.NewServeMux()
	combineMux.Handle(
		"/api/user/",
		http.StripPrefix("/api/user", middlewares.WithAuthorization(userHandlers)),
	)
	combineMux.Handle(
		"/api/auth/",
		http.StripPrefix("/api/auth", authHandlers),
	)
	combineMux.Handle(
		"/api/manage/",
		http.StripPrefix("/api/manage", middlewares.WithAuthorization(managingHandlers)),
	)
	combineMux.Handle(
		"/api/public/",
		http.StripPrefix("/api/public", publicHandlers),
	)
	combineMux.Handle(
		"/api/root/",
		http.StripPrefix("/api/root", middlewares.WithAuthorization(rootHandlers)),
	)

	serveError := http.Serve(listener, middlewares.WithLogger(combineMux))
	if serveError != nil {
		log.Fatal(serveError)
	}
}
