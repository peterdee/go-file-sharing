package manage

import (
	"net/http"

	"github.com/julyskies/gohelpers"

	"file-sharing/database"
	"file-sharing/utilities"
)

func DeleteFileHandler(response http.ResponseWriter, request *http.Request) {
	uid := request.PathValue("id")

	timestamp := gohelpers.MakeTimestampSeconds()
	database.FileService.UpdateOne(
		request.Context(),
		map[string]any{
			"isDeleted": false,
			"uid":       uid,
		},
		map[string]any{
			"$set": map[string]any{
				"deletedAt": timestamp,
				"isDeleted": true,
				"updatedAt": timestamp,
			},
		},
	)
	database.MetricsService.UpdateOne(
		request.Context(),
		map[string]any{
			"isDeleted": false,
			"uid":       uid,
		},
		map[string]any{
			"$set": map[string]any{
				"deletedAt": timestamp,
				"isDeleted": true,
				"updatedAt": timestamp,
			},
		},
	)

	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
