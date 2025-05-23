package public

import (
	"net/http"
	"os"
	"strconv"

	"github.com/julyskies/gohelpers"
	"github.com/nrednav/cuid2"
	"go.mongodb.org/mongo-driver/v2/bson"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
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
		DeletedAt:    0,
		IsDeleted:    false,
		OriginalName: handler.Filename,
		Size:         handler.Size,
		UID:          uid,
	}
	metricsRecord := database.Metrics{
		CreatedAt:      timestamp,
		DeletedAt:      0,
		Downloads:      0,
		IsDeleted:      false,
		LastDownloaded: timestamp,
		LastViewed:     timestamp,
		UID:            uid,
		Views:          0,
	}

	path := utilities.CreateFilePath(uid)
	destination, fileError := os.Create(path)
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

	insertError := database.Operations.InsertFile(filesRecord, request.Context())
	if insertError != nil {
		os.Remove(path)
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}
	insertError = database.Operations.InsertMetrics(metricsRecord, request.Context())
	if insertError != nil {
		database.Operations.DeleteFile(bson.M{"uid": uid}, request.Context())
		os.Remove(path)
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
