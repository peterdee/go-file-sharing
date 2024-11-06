package handlers

import (
	"net/http"

	"file-exchange/utilities"
)

func DownloadHandler(response http.ResponseWriter, request *http.Request) {
	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
