package structs

import (
	"encoding/json"
	"net/http"
)

type APIError struct {
	ErrorMessage string `json:"error"`
}

func (e APIError) With(detail string) APIError {
	e.ErrorMessage = detail
	return e
}

var Errors = map[string]APIError{
	"invalid_request":         {ErrorMessage: "invalid_request"},
	"unsupported_grant_type": {ErrorMessage: "unsupported_grant_type"},
	"server_error":           {ErrorMessage: "internal_server_error"},
}


func SendDetailedError(w http.ResponseWriter, err APIError, code int) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    json.NewEncoder(w).Encode(err)
}