package root

import (
	"net/http"

	"github.com/julyskies/gohelpers"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func DeleteUserHandler(response http.ResponseWriter, request *http.Request) {
	userData := middlewares.GetUserDataFromRequestContext(request.Context())
	if userData.Role != constants.ROLES.Root {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.Unauthorized,
			Request:  request,
			Response: response,
			Status:   http.StatusUnauthorized,
		})
		return
	}

	uid := request.PathValue("id")
	if userData.Uid == uid {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.Forbidden,
			InfoDetails: "Cannot delete own account",
			Request:     request,
			Response:    response,
			Status:      http.StatusForbidden,
		})
		return
	}

	cache.UserService.Del(request.Context(), uid)
	timestamp := gohelpers.MakeTimestampSeconds()
	database.UserService.UpdateOne(
		request.Context(),
		map[string]any{
			"isDeleted":      false,
			"setUpCompleted": true,
			"uid":            uid,
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
