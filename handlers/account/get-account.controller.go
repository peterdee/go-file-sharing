package account

import (
	"errors"
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func GetAccountHandler(response http.ResponseWriter, request *http.Request) {
	uid := middlewares.GetUidFromRequestContext(request.Context())
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
	queryError := database.UsersCollection.FindOne(
		request.Context(),
		bson.M{"uid": uid},
	).Decode(&user)
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

	utilities.Response(utilities.ResponseParams{
		Data: map[string]any{
			"user": user,
		},
		Request:  request,
		Response: response,
	})
}
