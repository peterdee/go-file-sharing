package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julyskies/gohelpers"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func ChangePasswordHandler(response http.ResponseWriter, request *http.Request) {
	parsed, parsingError := utilities.BodyParser(request, ChangePasswordRequestPayload{})
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

	newPassword := strings.Trim(parsed["newPassword"], " ")
	if newPassword == "" {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.BadRequest,
			InfoDetails: "Missing required 'newPassword' field",
			Request:     request,
			Response:    response,
			Status:      http.StatusBadRequest,
		})
		return
	}
	oldPassword := strings.Trim(parsed["oldPassword"], " ")
	if oldPassword == "" {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.BadRequest,
			InfoDetails: "Missing required 'oldPassword' field",
			Request:     request,
			Response:    response,
			Status:      http.StatusBadRequest,
		})
		return
	}

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
	queryError := database.Operations.GetUser(bson.M{"uid": uid}, &user, request.Context())
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

	match, hashError := utilities.ComparePlaintextWithHash(oldPassword, user.PasswordHash)
	if hashError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}
	if !match {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.Unauthorized,
			Request:  request,
			Response: response,
			Status:   http.StatusUnauthorized,
		})
		return
	}

	newPasswordHash, hashError := utilities.CreateHash(newPassword)
	if hashError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.InternalServerError,
			Request:  request,
			Response: response,
			Status:   http.StatusInternalServerError,
		})
		return
	}

	queryError = database.Operations.UpdateUser(
		bson.M{"uid": uid},
		bson.M{
			"$set": bson.M{
				"passwordHash": newPasswordHash,
				"updatedAt":    gohelpers.MakeTimestampSeconds(),
			},
		},
		request.Context(),
	)
	if queryError != nil {
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
