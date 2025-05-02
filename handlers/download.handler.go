package handlers

import (
	"fmt"
	"net/http"

	"file-exchange/utilities"
)

func DownloadHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Println("download handler")

	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
