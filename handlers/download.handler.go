package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

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

	uploadsDirectoryName := utilities.GetEnv(
		constants.ENV_NAMES.UplaodsDirectoryName,
		constants.DEFAULT_UPLOADS_DIRECTORY_NAME,
	)
	path := filepath.Join(uploadsDirectoryName, id)

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
		encoded, encodeError := json.Marshal(&filesRecord)
		if encodeError == nil {
			cache.Client.Set(
				context.Background(),
				id,
				encoded,
				time.Duration(time.Hour)*8,
			)
		}
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
	queryError := database.MetricsCollection.FindOneAndUpdate(
		context.Background(),
		bson.M{"uid": id},
		bson.M{
			"$inc": bson.M{"downloads": 1},
			"$set": bson.M{"lastDownloaded": gohelpers.MakeTimestampSeconds()},
		},
	).Decode(&metricsRecord)
	if queryError != nil {
		if errors.Is(queryError, mongo.ErrNoDocuments) {
			cache.Client.Del(context.Background(), id)
			database.FilesCollection.DeleteOne(context.Background(), bson.M{"uid": id})
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

	file, fileError := os.Open(path)
	if fileError != nil {
		if errors.Is(fileError, os.ErrNotExist) {
			cache.Client.Del(context.Background(), id)
			database.FilesCollection.DeleteOne(context.Background(), bson.M{"uid": id})
			database.MetricsCollection.DeleteOne(context.Background(), bson.M{"uid": id})
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
