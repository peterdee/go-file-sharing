package public

import (
	"net/http"

	"file-sharing/utilities"
)

func IndexHandler(response http.ResponseWriter, request *http.Request) {
	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
