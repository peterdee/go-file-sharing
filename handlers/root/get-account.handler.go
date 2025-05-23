package root

import (
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func GetUserHandler(response http.ResponseWriter, request *http.Request) {
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

	var user database.Users
	cacheError := cache.Operations.GetUser(uid, &user, request.Context())
	if cacheError != nil {
		cache.Operations.RemoveUser(uid, request.Context())
		queryError := database.Operations.GetUser(bson.M{"uid": uid}, &user, request.Context())
		if queryError != nil {
			if errors.Is(queryError, mongo.ErrNoDocuments) {
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
		cache.Operations.SaveUser(uid, user, request.Context())
	}

	utilities.Response(utilities.ResponseParams{
		Data:     map[string]any{"user": user},
		Request:  request,
		Response: response,
	})
}
