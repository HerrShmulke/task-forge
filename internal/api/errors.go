package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}

func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	errorResponse := ErrorResponse{
		Status: "error",
		Error:  message,
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(errorResponse)
}
