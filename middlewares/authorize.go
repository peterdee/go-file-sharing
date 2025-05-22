package middlewares

import (
	"context"
	"errors"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

type contextUserDataType string

const contextUserDataKey contextUserDataType = "contextUserDataKey"

type ContextUserData struct {
	Role string
	Uid  string
}

func GetUserDataFromRequestContext(requestContext context.Context) ContextUserData {
	return requestContext.Value(contextUserDataKey).(ContextUserData)
}

type Authorize struct {
	handler http.Handler
}

func (auth *Authorize) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	token := request.Header.Get("authorization")
	if token == "" {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.Unauthorized,
			InfoDetails: "Missing JWT",
			Request:     request,
			Response:    response,
			Status:      http.StatusUnauthorized,
		})
		return
	}
	uid, tokenError := utilities.ValidateJwt(token)
	if tokenError != nil {
		infoDetails := ""
		if errors.Is(tokenError, jwt.ErrSignatureInvalid) {
			infoDetails = "Token is invalid"
		}
		if errors.Is(tokenError, jwt.ErrTokenExpired) {
			infoDetails = "Token is expired"
		}
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.Unauthorized,
			InfoDetails: infoDetails,
			Request:     request,
			Response:    response,
			Status:      http.StatusUnauthorized,
		})
		return
	}

	// make sure that user exists
	var user database.Users
	cacheError := cache.Operations.GetUser(uid, user, request.Context())
	if cacheError != nil {
		cache.Operations.RemoveUser(uid, request.Context())
		queryError := database.Operations.GetUser(
			bson.M{
				"isDeleted":      false,
				"setUpCompleted": true,
				"uid":            uid,
			},
			user,
			request.Context(),
		)
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
		cache.Operations.SaveUser(uid, user, request.Context())
	}

	auth.handler.ServeHTTP(
		response,
		request.WithContext(
			context.WithValue(
				request.Context(),
				contextUserDataKey,
				ContextUserData{
					Role: user.Role,
					Uid:  uid,
				},
			),
		),
	)
}

func WithAuthorization(handler http.Handler) *Authorize {
	return &Authorize{handler}
}
