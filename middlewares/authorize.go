package middlewares

import (
	"context"
	"net/http"

	"file-sharing/constants"
	"file-sharing/utilities"
)

type Authorize struct {
	handler http.Handler
}

type contextUid string

const contextUidKey contextUid = "contextUidKey"

func GetUidFromRequestContext(requestContext context.Context) string {
	return requestContext.Value(contextUidKey).(string)
}

func (auth *Authorize) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	token := request.Header.Get("authorization")
	if token == "" {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.Unauthorized,
			InfoDetails: "Missing JWT",
			Request:     request,
			Response:    response,
			Status:      http.StatusUnauthorized,
		})
		return
	}
	uid, tokenError := utilities.ValidateJwt(token)
	if tokenError != nil {
		utilities.Response(utilities.ResponseParams{
			Info:        constants.RESPONSE_INFO.Unauthorized,
			InfoDetails: "Missing JWT",
			Request:     request,
			Response:    response,
			Status:      http.StatusUnauthorized,
		})
		return
	}

	auth.handler.ServeHTTP(
		response,
		request.WithContext(
			context.WithValue(
				request.Context(),
				contextUidKey,
				uid,
			),
		),
	)
}

func WithAuthorization(handler http.Handler) *Authorize {
	return &Authorize{handler}
}
