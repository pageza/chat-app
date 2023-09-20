// Package common provides common utilities and types for the chat application.
// It includes functionalities for handling API errors and sending error responses.
package common

import (
	"encoding/json"
	"net/http"
)

// APIError represents an error that can be sent in an API response.
type APIError struct {
	Status  int    `json:"status"`  // HTTP status code
	Message string `json:"message"` // Error message
}

// NewAPIError creates a new APIError with the given status code and message.
// It returns a pointer to the newly created APIError.
func NewAPIError(status int, message string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
	}
}

// RespondWithError sends an API response with an error message and status code.
// It sets the Content-Type header to "application/json" and writes the error to the response.
func RespondWithError(w http.ResponseWriter, err *APIError) {
	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the HTTP status code from the APIError
	w.WriteHeader(err.Status)

	// Encode the APIError as JSON and send it in the response
	json.NewEncoder(w).Encode(err)
}
