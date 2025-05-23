package utilities

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julyskies/gohelpers"

	"file-sharing/constants"
)

type ResponseParams struct {
	Data        any
	Info        string
	InfoDetails string
	Request     *http.Request
	Response    http.ResponseWriter
	Status      int
}

type responseObject struct {
	Data        any    `json:"data,omitempty"`
	Datetime    int64  `json:"datetime"`
	Info        string `json:"info"`
	InfoDetails string `json:"infoDetails,omitempty"`
	Request     string `json:"request"`
	Status      int    `json:"status"`
}

func Response(params ResponseParams) {
	info := params.Info
	if info == "" {
		info = constants.RESPONSE_INFO.Ok
	}

	status := params.Status
	if status == 0 {
		status = http.StatusOK
	}

	responseObject := responseObject{
		Datetime:    gohelpers.MakeTimestampSeconds(),
		Info:        info,
		InfoDetails: params.InfoDetails,
		Request:     fmt.Sprintf("%s [%s]", params.Request.RequestURI, params.Request.Method),
		Status:      status,
	}
	if params.Data != nil {
		responseObject.Data = params.Data
	}

	json, err := json.Marshal(responseObject)
	if err != nil {
		params.Response.WriteHeader(http.StatusInternalServerError)
		params.Response.Write([]byte(constants.RESPONSE_INFO.InternalServerError))
		return
	}

	params.Response.Header().Set("Content-Type", "application/json")
	params.Response.WriteHeader(status)
	params.Response.Write(json)
}
