package auth

import (
	"errors"
	"net/http"
	"strings"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/utilities"
)

func SignInHandler(response http.ResponseWriter, request *http.Request) {
	parsed, parsingError := utilities.BodyParser(request, SignInRequestPayload{})
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

	var user database.UserModel
	queryError := database.UserService.FindOne(
		request.Context(),
		map[string]any{
			"email":          email,
			"isDeleted":      false,
			"setUpCompleted": true,
		},
		&user,
	)
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

	match, hashError := utilities.ComparePlaintextWithHash(
		password,
		user.PasswordHash,
	)
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

	token, tokenError := utilities.CreateJwt(user.Uid)
	if tokenError != nil {
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
