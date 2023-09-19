// errors.go

package errors

import (
	"encoding/json"
	"net/http"
)

// NewAPIError creates a new APIError
func NewAPIError(status int, message string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
	}
}

// RespondWithError sends an API response with an error message and status code
func RespondWithError(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

// RespondWithCustomError sends an API response with a custom error type and status code
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
