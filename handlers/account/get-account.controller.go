package account

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func GetAccountHandler(response http.ResponseWriter, request *http.Request) {
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

	var user database.Users
	cachedUser, cacheError := getUserFromCache(uid, request.Context())
	if cacheError != nil {
		queryError := getUserFromDatabase(uid, &user, request.Context())
		if queryError != nil {
			if errors.Is(queryError, mongo.ErrNoDocuments) {
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
		saveUserToCache(uid, user, request.Context())
	}
	jsonError := json.Unmarshal([]byte(cachedUser), &user)
	if jsonError != nil {
		removeUserFromCache(uid, request.Context())
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}

	utilities.Response(utilities.ResponseParams{
		Data: map[string]any{
			"user": user,
		},
		Request:  request,
		Response: response,
	})
}
