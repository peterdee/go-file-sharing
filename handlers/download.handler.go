package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func DownloadHandler(response http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	cachedRecord, cacheError := getCachedRecord(id)
	if cacheError != nil || cachedRecord == nil {
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
		cacheError = setCacheValue(record.UID, &record)
		if cacheError != nil {
			utilities.Response(utilities.ResponseParams{
				Info:     constants.RESPONSE_INFO.InternalServerError,
				Request:  request,
				Response: response,
				Status:   http.StatusInternalServerError,
			})
			return
		}
		cachedRecord = &record
	}

	uploadsDirectoryName := utilities.GetEnv(
		constants.ENV_NAMES.UplaodsDirectoryName,
		constants.DEFAULT_UPLOADS_DIRECTORY_NAME,
	)
	file, fileError := os.Open(filepath.Join(uploadsDirectoryName, cachedRecord.UID))
	if fileError != nil {
		if errors.Is(fileError, os.ErrNotExist) {
			deleteCachedRecord(cachedRecord.UID)
			database.FilesCollection.DeleteOne(
				context.Background(),
				bson.D{{Key: "uid", Value: id}},
			)
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
		fmt.Sprintf("attachment; filename=%s", cachedRecord.OriginalName),
	)
	http.ServeFile(response, request, filepath.Join(uploadsDirectoryName, cachedRecord.UID))
}
