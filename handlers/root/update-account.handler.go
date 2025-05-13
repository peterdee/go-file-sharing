package root

import (
	"net/http"

	"file-sharing/utilities"
)

func UpdateAccountHandler(response http.ResponseWriter, request *http.Request) {
	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
