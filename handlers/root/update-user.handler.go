package root

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julyskies/gohelpers"

	"file-sharing/cache"
	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func UpdateUserHandler(response http.ResponseWriter, request *http.Request) {
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

	// check if user exists
	uid := request.PathValue("id")
	var user database.UserModel
	cacheError := cache.UserService.Get(request.Context(), uid, &user)
	if cacheError != nil {
		queryError := database.UserService.FindOneByUid(request.Context(), uid, &user)
		if queryError != nil {
			if errors.Is(queryError, database.ErrNoDocuments) {
				utilities.Response(utilities.ResponseParams{
					Info:        constants.RESPONSE_INFO.NotFound,
					InfoDetails: "User not found",
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

	update := map[string]any{}

	// check email
	email := strings.ToLower(strings.Trim(parsed["email"], " "))
	if email != user.Email {
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
			} else {
				update["email"] = email
			}
		}
	}

	// check password
	password := strings.Trim(parsed["password"], " ")
	if password != "" {
		hashed, hashError := utilities.CreateHash(password)
		if hashError != nil {
			utilities.Response(utilities.ResponseParams{
				Info:     constants.RESPONSE_INFO.InternalServerError,
				Request:  request,
				Response: response,
				Status:   http.StatusInternalServerError,
			})
			return
		}
		update["passwordHash"] = hashed
	}

	// make sure that role is correct
	role := strings.ToLower(strings.Trim(parsed["role"], " "))
	if role != "" && !gohelpers.IncludesString(gohelpers.StructValues(constants.ROLES), role) {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.BadRequest,
			InfoDetails: "Invalid role",
			Request:     request,
			Response:    response,
			Status:      http.StatusBadRequest,
		})
		return
	}
	if role != user.Role {
		update["role"] = role
	}

	// apply updates only if there are anything to be updated
	if len(update) > 0 {
		update["updatedAt"] = gohelpers.MakeTimestampSeconds()
		queryError := database.UserService.UpdateOne(
			request.Context(),
			map[string]any{"uid": uid},
			update,
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
		cache.UserService.Del(request.Context(), uid)
	}

	utilities.Response(utilities.ResponseParams{
		Request:  request,
		Response: response,
	})
}
