package public

import (
	"errors"
	"net/http"
	"os"

	"github.com/julyskies/gohelpers"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func InfoHandler(response http.ResponseWriter, request *http.Request) {
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

	var metrics database.MetricsModel
	timestamp := gohelpers.MakeTimestampSeconds()
	queryError := database.MetricsService.FindOneAndUpdate(
		request.Context(),
		map[string]any{"uid": uid},
		map[string]any{
			"$inc": map[string]any{"views": 1},
			"$set": map[string]any{
				"lastViewed": timestamp,
				"updatedAt":  timestamp,
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
