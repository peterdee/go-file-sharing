package root

import (
	"net/http"

	"file-sharing/constants"
	"file-sharing/database"
	"file-sharing/middlewares"
	"file-sharing/utilities"
)

func ListUsersHandler(response http.ResponseWriter, request *http.Request) {
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

	pagination := utilities.RequestPagination(request)
	var users []database.UserModel
	count, queryError := database.UserService.FindPaginated(
		request.Context(),
		map[string]any{},
		pagination,
		&users,
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

	// send an empty array instead of nil
	if count == 0 || users == nil {
		users = []database.UserModel{}
	}

	utilities.Response(utilities.ResponseParams{
		Data: map[string]any{
			"pagination": utilities.ResponsePagination(pagination, count),
			"users":      users,
		},
		Request:  request,
		Response: response,
	})
}
