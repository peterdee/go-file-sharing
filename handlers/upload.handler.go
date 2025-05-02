package handlers

import (
	"fmt"
	"net/http"

	"file-exchange/utilities"
)

func UploadHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("upload handler")

	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
