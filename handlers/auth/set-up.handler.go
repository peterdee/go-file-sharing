package auth

import (
	"errors"
	"net/http"

	"github.com/julyskies/gohelpers"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func SetUpHandler(response http.ResponseWriter, request *http.Request) {
	parsed, parsingError := utilities.BodyParser(request, SetUpRequestPayload{})
	if parsingError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.BadRequest,
			InfoDetails: parsingError.Error(),
			Request:     request,
			Response:    response,
			Status:      http.StatusBadRequest,
		})
		return
	}

	email := parsed["email"]
	if email == "" {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.BadRequest,
			InfoDetails: "Missing required 'email' field",
			Request:     request,
			Response:    response,
			Status:      http.StatusBadRequest,
		})
		return
	}
	password := parsed["password"]
	if password == "" {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.BadRequest,
			InfoDetails: "Missing required 'password' field",
			Request:     request,
			Response:    response,
			Status:      http.StatusBadRequest,
		})
		return
	}

	passwordHash, hashError := utilities.CreateHash(password)
	if hashError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}

	var user database.Users
	queryError := database.UsersCollection.FindOneAndUpdate(
		request.Context(),
		bson.M{
			"email":          email,
			"setUpCompleted": false,
		},
		bson.M{
			"passwordHash":   passwordHash,
			"setUpCompleted": true,
			"updatedAt":      gohelpers.MakeTimestampSeconds(),
		},
	).Decode(&user)
	if queryError != nil {
		if errors.Is(queryError, mongo.ErrNoDocuments) {
			utilities.Response(utilities.ResponseParams{
				Info:        constants.RESPONSE_INFO.NotFound,
				InfoDetails: "Account not found",
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

	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
