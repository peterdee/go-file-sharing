package user

import (
	"errors"
	"fmt"
	"net/http"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func GetUserHandler(response http.ResponseWriter, request *http.Request) {
	userData := middlewares.GetUserDataFromRequestContext(request.Context())
	uid := userData.Uid
	if uid == "" {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}

	var user database.UserModel
	cacheError := cache.UserService.Get(request.Context(), uid, &user)
	if cacheError != nil {
		queryError := database.UserService.FindOneByUid(request.Context(), uid, &user)
		fmt.Println(queryError)
		if queryError != nil {
			if errors.Is(queryError, database.ErrNoDocuments) {
				utilities.Response(utilities.ResponseParams{
					Info:     constants.RESPONSE_INFO.Unauthorized,
					Request:  request,
					Response: response,
					Status:   http.StatusUnauthorized,
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
		fmt.Println(user)
		cache.UserService.Set(request.Context(), user)
	}

	utilities.Response(utilities.ResponseParams{
		Data:     map[string]any{"user": user},
		Request:  request,
		Response: response,
	})
}
