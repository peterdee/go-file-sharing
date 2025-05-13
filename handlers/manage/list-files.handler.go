package manage

import (
	"net/http"

	"file-sharing/utilities"
)

func ListFilesHandler(response http.ResponseWriter, request *http.Request) {
	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
