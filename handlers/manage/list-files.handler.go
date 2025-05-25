package manage

import (
	"fmt"
	"net/http"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func ListFilesHandler(response http.ResponseWriter, request *http.Request) {
	pagination := utilities.RequestPagination(request)
	var files []database.FileModel
	count, queryError := database.FileService.FindPaginated(
		request.Context(),
		map[string]any{},
		pagination,
		&files,
	)
	if queryError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}

	if count == 0 || files == nil {
		utilities.Response(utilities.ResponseParams{
			Data: map[string]any{
				"files":      []database.FileModel{},
				"pagination": utilities.ResponsePagination(pagination, count),
			},
			Request:  request,
			Response: response,
		})
		return
	}

	fileUids := make([]string, len(files))
	for index, file := range files {
		fileUids[index] = file.Uid
	}

	var metrics []database.MetricsModel
	queryError = database.MetricsService.FindAll(
		request.Context(),
		map[string]any{"uid": map[string]any{"$in": fileUids}},
		&metrics,
	)
	if queryError != nil {
		fmt.Println(queryError)
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}

	filesWithStats := make([]map[string]any, len(files))
	for index, file := range files {
		for i := 0; i < len(metrics); i += 1 {
			if metrics[i].Uid == file.Uid {
				filesWithStats[index] = map[string]any{
					"file":  file,
					"stats": metrics[i],
				}
				break
			}
		}
	}

	utilities.Response(utilities.ResponseParams{
		Data: map[string]any{
			"files":      filesWithStats,
			"pagination": utilities.ResponsePagination(pagination, count),
		},
		Request:  request,
		Response: response,
	})
}
