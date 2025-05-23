package user

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julyskies/gohelpers"

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

	var user database.UserModel
	queryError := database.UserService.FindOneByUid(request.Context(), uid, &user)
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

	queryError = database.UserService.UpdateOne(
		request.Context(),
		map[string]any{"uid": uid},
		map[string]any{
			"$set": map[string]any{
				"passwordHash": newPasswordHash,
				"updatedAt":    gohelpers.MakeTimestampSeconds(),
			},
		},
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
