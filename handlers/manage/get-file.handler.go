package manage

import (
	"errors"
	"net/http"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func GetFileHandler(response http.ResponseWriter, request *http.Request) {
	uid := request.PathValue("id")

	var file database.FileModel
	cacheError := cache.FileService.Get(request.Context(), uid, &file)
	if cacheError != nil {
		queryError := database.FileService.FindOneByUid(request.Context(), uid, &file)
		if queryError != nil {
			if errors.Is(queryError, database.ErrNoDocuments) {
				utilities.Response(utilities.ResponseParams{
					Info:        constants.RESPONSE_INFO.NotFound,
					InfoDetails: "File not found",
					Request:     request,
					Response:    response,
					Status:      http.StatusNotFound,
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
	}

	var metrics database.MetricsModel
	queryError := database.MetricsService.FindOneByUid(request.Context(), uid, &metrics)
	if queryError != nil {
		if errors.Is(queryError, database.ErrNoDocuments) {
			utilities.Response(utilities.ResponseParams{
				Info:        constants.RESPONSE_INFO.NotFound,
				InfoDetails: "File not found",
				Request:     request,
				Response:    response,
				Status:      http.StatusNotFound,
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

	utilities.Response(utilities.ResponseParams{
		Data: map[string]any{
			"file":  file,
			"stats": metrics,
		},
		Request:  request,
		Response: response,
	})
}
