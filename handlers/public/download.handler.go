package public

import (
	"encoding/json"
	"errors"
	"fmt"
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

func DownloadHandler(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")
	path := createFilePath(id)

	var filesRecord database.Files

	cachedFilesRecord, cacheError := cache.Client.Get(
		request.Context(),
		id,
	).Result()
	if cacheError == nil {
		cacheError = json.Unmarshal([]byte(cachedFilesRecord), &filesRecord)
	}

	if cacheError != nil {
		queryError := database.FilesCollection.FindOne(
			request.Context(),
			bson.M{"uid": id},
		).Decode(&filesRecord)
		if queryError != nil {
			if errors.Is(queryError, mongo.ErrNoDocuments) {
				database.MetricsCollection.DeleteOne(request.Context(), bson.M{"uid": id})
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

	if filesRecord.IsDeleted {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.FileNotAvailable,
			Request:  request,
			Response: response,
			Status:   http.StatusGone,
		})
		return
	}

	var metricsRecord database.Metrics
	timestamp := gohelpers.MakeTimestampSeconds()
	queryError := database.MetricsCollection.FindOneAndUpdate(
		request.Context(),
		bson.M{"uid": id},
		bson.M{
			"$inc": bson.M{"downloads": 1},
			"$set": bson.M{
				"lastDownloaded": timestamp,
				"updatedAt":      timestamp,
			},
		},
	).Decode(&metricsRecord)
	if queryError != nil {
		if errors.Is(queryError, mongo.ErrNoDocuments) {
			database.FilesCollection.DeleteOne(request.Context(), bson.M{"uid": id})
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

	file, fileError := os.Open(path)
	if fileError != nil {
		if errors.Is(fileError, os.ErrNotExist) {
			database.FilesCollection.DeleteOne(request.Context(), bson.M{"uid": id})
			database.MetricsCollection.DeleteOne(request.Context(), bson.M{"uid": id})
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
	defer file.Close()

	response.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", filesRecord.OriginalName),
	)
	http.ServeFile(response, request, path)
}
