package utilities

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

type BodyParserError struct {
	message string
}

func (instance *BodyParserError) Error() string {
	return instance.message
}

func BodyParser(
	response http.ResponseWriter,
	request *http.Request,
	destination any,
) error {
	// contentType := request.Header.Get("Content-Type")
	// if contentType != "" {
	// 	mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
	// 	if mediaType != "application/json" {
	// 		return &BodyParserError{"content-type header is not application/json"}
	// 	}
	// }

	request.Body = http.MaxBytesReader(response, request.Body, 1048576) // 1 MB

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	decodeError := decoder.Decode(&destination)
	if decodeError != nil {
		if errors.Is(decodeError, io.EOF) {
			return &BodyParserError{"request body is empty"}
		}
		if strings.HasPrefix(decodeError.Error(), "json: unknown field") {
			return &BodyParserError{"request body contains unknown field"}
		}
		if errors.Is(decodeError, io.ErrUnexpectedEOF) {
			return &BodyParserError{"request body contains malformed JSON"}
		}
		if errors.Is(decodeError, &json.UnmarshalTypeError{}) {
			return &BodyParserError{"request body contains an invalid value for the field"}
		}
		return decodeError
	}

	return nil
}
