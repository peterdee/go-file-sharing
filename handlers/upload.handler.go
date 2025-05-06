package handlers

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

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
	defer file.Close()

	uid := cuid2.Generate()
	fileRecord := database.File{
		CreatedAt:    gohelpers.MakeTimestampSeconds(),
		Downloads:    0,
		OriginalName: handler.Filename,
		Size:         handler.Size,
		UID:          uid,
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
		fileRecord,
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
