package handlers

import (
	"context"
	"net/http"

	"file-exchange/database"
	"file-exchange/utilities"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func InfoHandler(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	result := database.FilesCollection.FindOne(
		context.Background(),
		bson.D{{Key: "UID", Value: id}},
	)

	/*
		1. Get database record by ID
		2. Return 404 error if record was not found
		3. Return response
	*/

	utilities.Response(utilities.ResponseParams{
		Data:     result,
		Request:  request,
		Response: response,
	})
}
