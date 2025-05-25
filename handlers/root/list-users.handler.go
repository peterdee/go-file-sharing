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

	pagination := utilities.Pagination(request)
	var users []database.UserModel
	queyError := database.UserService.FindPaginated(request.Context(), pagination, &users)
	if queyError != nil {
		// TODO: handle
	}

	utilities.Response(utilities.ResponseParams{
		Data: map[string]any{
			"users": users,
		},
		Request:  request,
		Response: response,
	})
}
