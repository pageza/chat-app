package common

import (
	"encoding/json"
	"net/http"
)

// APIError represents an error that can be sent in an API response
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// NewAPIError creates a new APIError
func NewAPIError(status int, message string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
	}
}

// RespondWithError sends an API response with a error message and status code
func RespondWithError(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}
