package utilities

import (
	"math"
	"net/http"
	"strconv"
)

type PaginationData struct {
	Limit  int
	Offset int
	Page   int
}

func getFromQuery(request *http.Request, key string, defaultValue int) int {
	stringValue := request.URL.Query().Get(key)
	if stringValue != "" {
		value, typeError := strconv.Atoi(stringValue)
		if typeError == nil {
			return value
		}
		return defaultValue
	}
	return defaultValue
}

func RequestPagination(request *http.Request) PaginationData {
	limit := getFromQuery(request, "limit", 10)
	page := getFromQuery(request, "page", 1)

	return PaginationData{
		Limit:  limit,
		Offset: (page - 1) * limit,
		Page:   page,
	}
}

func ResponsePagination(pagination PaginationData, totalCount int64) map[string]any {
	return map[string]any{
		"limit":      pagination.Limit,
		"page":       pagination.Page,
		"totalCount": totalCount,
		"totalPages": math.Ceil(float64(totalCount) / float64(pagination.Limit)),
	}
}
