package handlers

import (
	"fmt"
	"net/http"

	"file-exchange/utilities"
)

func InfoHandler(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	fmt.Println("info handler", id)

	/*
		1. Get database record by ID
		2. Return 404 error if record was not found
		3. Return response
	*/

	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
