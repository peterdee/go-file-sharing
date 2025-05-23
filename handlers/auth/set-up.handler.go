package auth

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julyskies/gohelpers"

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

	var user database.UserModel
	queryError := database.UserService.FindOneAndUpdate(
		request.Context(),
		map[string]any{
			"email":          email,
			"isDeleted":      false,
			"setUpCompleted": false,
		},
		map[string]any{
			"$set": map[string]any{
				"passwordHash":   passwordHash,
				"setUpCompleted": true,
				"updatedAt":      gohelpers.MakeTimestampSeconds(),
			},
		},
		&user,
	)
	if queryError != nil {
		if errors.Is(queryError, database.ErrNoDocuments) {
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
