package handlers

import (
	"fmt"
	"net/http"

	"file-exchange/utilities"
)

func IndexHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("index handler")

	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
