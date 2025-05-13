package auth

import (
	"errors"
	"fmt"
	"net/http"

	"file-sharing/constants"
	"file-sharing/utilities"
)

func SetUpHandler(response http.ResponseWriter, request *http.Request) {
	var payload setUpRequestPayload
	decodeError := utilities.BodyParser(response, request, payload)

	if decodeError != nil {
		fmt.Println(decodeError)
		if errors.Is(decodeError, &utilities.BodyParserError{}) {
			utilities.Response(utilities.ResponseParams{
				Info:        constants.RESPONSE_INFO.BadRequest,
				InfoDetails: decodeError.Error(),
				Request:     request,
				Response:    response,
				Status:      http.StatusBadRequest,
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
