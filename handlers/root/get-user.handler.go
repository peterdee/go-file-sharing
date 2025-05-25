package root

import (
	"errors"
	"net/http"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func GetUserHandler(response http.ResponseWriter, request *http.Request) {
	authData := middlewares.GetUserDataFromRequestContext(request.Context())
	if authData.Role != constants.ROLES.Root {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.Unauthorized,
			Request:  request,
			Response: response,
			Status:   http.StatusUnauthorized,
		})
		return
	}

	uid := request.PathValue("id")

	var user database.UserModel
	cacheError := cache.UserService.Get(request.Context(), uid, &user)
	if cacheError != nil {
		queryError := database.UserService.FindOneByUid(request.Context(), uid, &user)
		if queryError != nil {
			if errors.Is(queryError, database.ErrNoDocuments) {
				utilities.Response(utilities.ResponseParams{
					Info:        constants.RESPONSE_INFO.NotFound,
					InfoDetails: "User account not found",
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
		cache.UserService.Set(request.Context(), user)
	}

	utilities.Response(utilities.ResponseParams{
		Data:     map[string]any{"user": user},
		Request:  request,
		Response: response,
	})
}
