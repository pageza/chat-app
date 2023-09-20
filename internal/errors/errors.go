// Package errors provides utility functions for creating and responding with custom error types.
// This package is designed to make error handling more structured and informative across the application.

package errors

import (
	"encoding/json"
	"net/http"
)

// NewAPIError creates a new APIError instance.
// This function is useful for generating API errors with a specific HTTP status and message.
//
// Parameters:
// - status: HTTP status code for the error
// - message: Error message to be displayed
//
// Returns:
// - A pointer to a new APIError instance
func NewAPIError(status int, message string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
	}
}

// RespondWithError sends an API response containing an APIError.
// This function sets the HTTP status code and Content-Type header before sending the error as JSON.
//
// Parameters:
// - w: The http.ResponseWriter to write the response to
// - err: The APIError to send in the response
func RespondWithError(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

// RespondWithCustomError sends an API response containing a custom error type.
// This function sets the HTTP status code based on the custom error type and sends the error as JSON.
//
// Parameters:
// - w: The http.ResponseWriter to write the response to
// - err: The custom error to send in the response
func RespondWithCustomError(w http.ResponseWriter, err interface{}) {
	w.Header().Set("Content-Type", "application/json")
	switch e := err.(type) {
	case *AuthenticationError:
		w.WriteHeader(e.Status)
	case *RateLimitError:
		w.WriteHeader(e.Status)
	case *DatabaseError:
		w.WriteHeader(e.Status)
	case *ValidationError:
		w.WriteHeader(e.Status)
	case *InternalServerError:
		w.WriteHeader(e.Status)
	case *NotFoundError:
		w.WriteHeader(e.Status)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(err)
}
