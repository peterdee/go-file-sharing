package public

import (
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/julyskies/gohelpers"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func DownloadHandler(response http.ResponseWriter, request *http.Request) {
	uid := request.PathValue("id")
	path := utilities.CreateFilePath(uid)

	var file database.FileModel
	cacheError := cache.FileService.Get(request.Context(), uid, &file)
	if cacheError != nil {
		queryError := database.FileService.FindOneByUid(request.Context(), uid, &file)
		if queryError != nil {
			if errors.Is(queryError, database.ErrNoDocuments) {
				database.MetricsService.DeleteOneByUid(request.Context(), uid)
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
		cache.FileService.Set(request.Context(), file)
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

	var metrics database.MetricsModel
	timestamp := gohelpers.MakeTimestampSeconds()
	queryError := database.MetricsService.FindOneAndUpdate(
		request.Context(),
		map[string]any{"uid": uid},
		map[string]any{
			"$inc": map[string]any{"downloads": 1},
			"$set": map[string]any{
				"lastDownloaded": timestamp,
				"updatedAt":      timestamp,
			},
		},
		&metrics,
	)
	if queryError != nil {
		if errors.Is(queryError, database.ErrNoDocuments) {
			cache.FileService.Del(request.Context(), uid)
			database.FileService.DeleteOneByUid(request.Context(), uid)
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

	fileData, fileError := os.Open(path)
	if fileError != nil {
		if errors.Is(fileError, os.ErrNotExist) {
			cache.FileService.Del(request.Context(), uid)
			database.FileService.DeleteOneByUid(request.Context(), uid)
			database.MetricsService.DeleteOneByUid(request.Context(), uid)
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
