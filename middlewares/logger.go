package middlewares

import (
	"log"
	"net/http"
	"time"
)

type Logger struct {
	handler http.Handler
}

func (logger *Logger) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	start := time.Now()
	logger.handler.ServeHTTP(response, request)
	log.Printf(
		"%s %s - %v",
		request.Method,
		request.URL.Path,
		time.Since(start),
	)
}

func WithLogger(handler http.Handler) *Logger {
	return &Logger{handler}
}
