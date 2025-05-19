package account

import (
	"net/http"

	"file-sharing/utilities"
)

func GetAccountHandler(response http.ResponseWriter, request *http.Request) {
	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
