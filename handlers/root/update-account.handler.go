package root

import (
	"net/http"

	"file-sharing/utilities"
)

func UpdateUserHandler(response http.ResponseWriter, request *http.Request) {
	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
