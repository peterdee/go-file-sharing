package public

import (
	"errors"
	"net/http"
	"os"

	"github.com/julyskies/gohelpers"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func InfoHandler(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")
	path := utilities.CreateFilePath(id)

	var file database.Files
	cacheError := cache.Operations.GetFile(id, file, request.Context())
	if cacheError != nil {
		queryError := database.Operations.GetFile(bson.M{"uid": id}, &file, request.Context())
		if queryError != nil {
			if errors.Is(queryError, mongo.ErrNoDocuments) {
				database.Operations.DeleteMetrics(bson.M{"uid": id}, request.Context())
				os.Remove(path)
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
		cache.Operations.SaveFile(id, file, request.Context())
	}

	var metrics database.Metrics
	timestamp := gohelpers.MakeTimestampSeconds()
	queryError := database.Operations.GetMetricsAndUpdate(
		bson.M{"uid": id},
		bson.M{
			"$inc": bson.M{"views": 1},
			"$set": bson.M{
				"lastViewed": timestamp,
				"updatedAt":  timestamp,
			},
		},
		&metrics,
		request.Context(),
	)
	if queryError != nil {
		if errors.Is(queryError, mongo.ErrNoDocuments) {
			database.Operations.DeleteFile(bson.M{"uid": id}, request.Context())
			os.Remove(path)
			cache.Operations.RemoveFile(id, request.Context())
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

	metrics.LastViewed = timestamp
	metrics.Views += 1

	utilities.Response(utilities.ResponseParams{
		Data: map[string]any{
			"file":  file,
			"stats": metrics,
		},
		Request:  request,
		Response: response,
	})
}
