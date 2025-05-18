package public

import (
	"context"
	"encoding/json"
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
	path := createFilePath(id)

	var filesRecord database.Files

	cachedFilesRecord, cacheError := cache.Client.Get(
		context.Background(),
		id,
	).Result()
	if cacheError == nil {
		cacheError = json.Unmarshal([]byte(cachedFilesRecord), &filesRecord)
	}

	if cacheError != nil {
		queryError := database.FilesCollection.FindOne(
			context.Background(),
			bson.M{"uid": id},
		).Decode(&filesRecord)
		if queryError != nil {
			if errors.Is(queryError, mongo.ErrNoDocuments) {
				database.MetricsCollection.DeleteOne(context.Background(), bson.M{"uid": id})
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
		saveToCache(id, filesRecord)
	}

	var metricsRecord database.Metrics
	timestamp := gohelpers.MakeTimestampSeconds()
	queryError := database.MetricsCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"uid": id},
		bson.M{
			"$inc": bson.M{"views": 1},
			"$set": bson.M{
				"lastViewed": timestamp,
				"updatedAt":  timestamp,
			},
		},
	).Decode(&metricsRecord)
	if queryError != nil {
		if errors.Is(queryError, mongo.ErrNoDocuments) {
			database.FilesCollection.DeleteOne(context.Background(), bson.M{"uid": id})
			os.Remove(path)
			removeFromCache(id)
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

	metricsRecord.LastViewed = timestamp
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
