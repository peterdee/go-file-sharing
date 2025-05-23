package public

import (
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

	if file.IsDeleted {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.FileNotAvailable,
			Request:  request,
			Response: response,
			Status:   http.StatusGone,
		})
		return
	}

	var metrics database.Metrics
	timestamp := gohelpers.MakeTimestampSeconds()
	queryError := database.Operations.GetMetricsAndUpdate(
		bson.M{"uid": id},
		bson.M{
			"$inc": bson.M{"downloads": 1},
			"$set": bson.M{
				"lastDownloaded": timestamp,
				"updatedAt":      timestamp,
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

	fileData, fileError := os.Open(path)
	if fileError != nil {
		if errors.Is(fileError, os.ErrNotExist) {
			database.Operations.DeleteFile(bson.M{"uid": id}, request.Context())
			database.Operations.DeleteMetrics(bson.M{"uid": id}, request.Context())
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
	defer fileData.Close()

	response.Header().Set(
		"Content-Disposition",
		fmt.Sprintf("attachment; filename=%s", file.OriginalName),
	)
	http.ServeFile(response, request, path)
}
