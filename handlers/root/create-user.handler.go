package root

import (
	"context"
	"net/http"
	"strings"

	"github.com/julyskies/gohelpers"
	"github.com/nrednav/cuid2"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func CreateUserHandler(response http.ResponseWriter, request *http.Request) {
	authData := middlewares.GetUserDataFromRequestContext(request.Context())
	if authData.Role != constants.ROLES.Root {
		utilities.Response(utilities.ResponseParams{
			Info:     constants.RESPONSE_INFO.Unauthorized,
			Request:  request,
			Response: response,
			Status:   http.StatusUnauthorized,
		})
		return
	}

	parsed, parsingError := utilities.BodyParser(request, CreateUserRequestPayload{})
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

	queryError := database.UserService.FindOne(
		request.Context(),
		map[string]any{"email": email},
		&database.UserModel{},
	)
	if queryError == nil {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.BadRequest,
			InfoDetails: "Email address is already in use",
			Request:     request,
			Response:    response,
			Status:      http.StatusBadRequest,
		})
		return
	} else {
		if queryError != database.ErrNoDocuments {
			utilities.Response(utilities.ResponseParams{
				Info:     constants.RESPONSE_INFO.InternalServerError,
				Request:  request,
				Response: response,
				Status:   http.StatusInternalServerError,
			})
			return
		}
	}

	role := strings.ToLower(strings.Trim(parsed["role"], " "))
	if role == "" {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.BadRequest,
			InfoDetails: "Missing required 'role' field",
			Request:     request,
			Response:    response,
			Status:      http.StatusBadRequest,
		})
		return
	}
	if !gohelpers.IncludesString(gohelpers.StructValues(constants.ROLES), role) {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.BadRequest,
			InfoDetails: "Invalid role",
			Request:     request,
			Response:    response,
			Status:      http.StatusBadRequest,
		})
		return
	}

	timestamp := gohelpers.MakeTimestampSeconds()
	user := database.UserModel{
		CreatedAt:      timestamp,
		DeletedAt:      0,
		Email:          email,
		IsDeleted:      false,
		PasswordHash:   "",
		Role:           role,
		SetUpCompleted: false,
		Uid:            cuid2.Generate(),
		UpdatedAt:      timestamp,
	}
	queryError = database.UserService.InsertOne(context.Background(), user)
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
		Data:     map[string]any{"user": user},
		Request:  request,
		Response: response,
	})
}
