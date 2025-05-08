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

	var filesRecord database.Files
	queryError := database.FilesCollection.FindOne(
		context.Background(),
		bson.M{"uid": id},
	).Decode(&filesRecord)
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

	var metricsRecord database.Metrics
	queryError = database.MetricsCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"uid": id},
		bson.M{"$inc": bson.M{"views": 1}},
	).Decode(&metricsRecord)
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

	metricsRecord.Views += 1

	utilities.Response(utilities.ResponseParams{
		Data: map[string]any{
			"file":  filesRecord,
			"stats": metricsRecord,
		},
		Request:  request,
		Response: response,
	})
}
