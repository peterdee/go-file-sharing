package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"

	"github.com/julyskies/gohelpers"
	"github.com/nrednav/cuid2"
)

func UploadHandler(response http.ResponseWriter, request *http.Request) {
	file, handler, fileError := request.FormFile("file")
	if fileError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}

	// no size restrictions if MAX_FILE_SIZE_BYTES is set to 0 or is not set at all
	maxFileSizeBytes := utilities.GetEnv(constants.ENV_NAMES.MaxFileSizeBytes)
	if maxFileSizeBytes != "" && maxFileSizeBytes != "0" {
		size, converterError := strconv.Atoi(maxFileSizeBytes)
		// if value is not an int consider that there are no restrictions
		if converterError == nil && handler.Size > int64(size) {
			utilities.Response(utilities.ResponseParams{
				Info:     constants.RESPONSE_INFO.RequestEntityTooLarge,
				Request:  request,
				Response: response,
				Status:   http.StatusRequestEntityTooLarge,
			})
			return
		}
	}

	defer file.Close()

	timestamp := gohelpers.MakeTimestampSeconds()
	uid := cuid2.Generate()
	filesRecord := database.Files{
		CreatedAt:    timestamp,
		OriginalName: handler.Filename,
		Size:         handler.Size,
		UID:          uid,
	}
	metricsRecord := database.Metrics{
		CreatedAt: timestamp,
		Downloads: 0,
		UID:       uid,
		Views:     0,
	}

	uploadsDirectoryName := utilities.GetEnv(
		constants.ENV_NAMES.UplaodsDirectoryName,
		constants.DEFAULT_UPLOADS_DIRECTORY_NAME,
	)
	destination, fileError := os.Create(filepath.Join(uploadsDirectoryName, uid))
	if fileError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}
	defer destination.Close()

	if _, copyError := destination.ReadFrom(file); copyError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}

	_, insertError := database.FilesCollection.InsertOne(
		context.Background(),
		filesRecord,
	)
	if insertError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}
	_, insertError = database.MetricsCollection.InsertOne(
		context.Background(),
		metricsRecord,
	)
	if insertError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}

	utilities.Response(utilities.ResponseParams{
		Data: map[string]string{
			"id": uid,
		},
		Request:  request,
		Response: response,
	})
}
