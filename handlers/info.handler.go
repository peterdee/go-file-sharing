package handlers

import (
	"context"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func InfoHandler(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	var record database.File
	queryError := database.FilesCollection.FindOne(
		context.Background(),
		bson.D{{Key: "uid", Value: id}},
	).Decode(&record)
	if queryError != nil {
		if errors.Is(queryError, mongo.ErrNoDocuments) {
			utilities.Response(utilities.ResponseParams{
				Info:     constants.RESPONSE_INFO.NotFound,
				Request:  request,
				Response: response,
				Status:   http.StatusNotFound,
			})
			return
		}
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}

	utilities.Response(utilities.ResponseParams{
		Data:     record,
		Request:  request,
		Response: response,
	})
}
