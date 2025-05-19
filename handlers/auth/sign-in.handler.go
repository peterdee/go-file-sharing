package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SignInHandler(response http.ResponseWriter, request *http.Request) {
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

	email := strings.ToLower(strings.Trim(parsed["email"], " "))
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
	password := strings.Trim(parsed["password"], " ")
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

	var user database.Users
	queryError := database.UsersCollection.FindOne(
		request.Context(),
		bson.M{
			"email":          email,
			"isDeleted":      false,
			"setUpCompleted": true,
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

	token, tokenError := utilities.CreateJwt(user.UID)
	if tokenError != nil {
		fmt.Println(tokenError)
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
			"token": token,
			"user":  user,
		},
		Request:  request,
		Response: response,
	})
}
