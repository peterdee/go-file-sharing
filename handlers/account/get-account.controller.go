package account

import (
	"fmt"
	"net/http"

	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func GetAccountHandler(response http.ResponseWriter, request *http.Request) {
	uid := middlewares.GetUidFromRequestContext(request.Context())

	fmt.Println(uid)

	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
