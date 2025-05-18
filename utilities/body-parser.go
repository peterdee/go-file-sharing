package utilities

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julyskies/gohelpers"
)

func BodyParser(
	request *http.Request,
	payloadStruct any,
) (map[string]string, error) {
	contentType := request.Header.Get("Content-Type")
	if contentType != "" {
		mediaType := strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
		if mediaType != "application/x-www-form-urlencoded" {
			return nil, errors.New("Content-Type header is not application/x-www-form-urlencoded")
		}
	}

	fieldNames, _ := gohelpers.StructFieldsJson(
		payloadStruct,
		gohelpers.StructKeysJsonParams{
			SkipIgnoredFields: true,
			SkipMissingFields: true,
		},
	)

	if parsingError := request.ParseForm(); parsingError != nil {
		return nil, errors.New("Malformed request payload")
	}

	result := make(map[string]string)
	for _, field := range fieldNames {
		result[field] = request.FormValue(field)
	}
	return result, nil
}
