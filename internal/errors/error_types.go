// Package errors defines custom error types that are used throughout the chat application.
// These error types are designed to provide detailed information about different kinds of errors
// that can occur, such as authentication errors, rate limit errors, and database errors.

package errors

// AuthenticationError represents errors related to user authentication.
// This could be due to invalid credentials, expired tokens, etc.
type AuthenticationError struct {
	Status  int    `json:"status"`           // HTTP status code
	Message string `json:"message"`          // Error message
	Reason  string `json:"reason,omitempty"` // Optional additional information about the error
}

// RateLimitError represents errors related to exceeding API rate limits.
type RateLimitError struct {
	Status     int    `json:"status"`      // HTTP status code
	Message    string `json:"message"`     // Error message
	RetryAfter int    `json:"retry_after"` // Time in seconds after which the client may retry
}

// DatabaseError represents errors related to database operations.
type DatabaseError struct {
	Status  int    `json:"status"`          // HTTP status code
	Message string `json:"message"`         // Error message
	Query   string `json:"query,omitempty"` // Optional SQL query that caused the error
}

// ValidationError represents errors related to invalid input or validation failure.
type ValidationError struct {
	Status  int      `json:"status"`           // HTTP status code
	Message string   `json:"message"`          // Error message
	Fields  []string `json:"fields,omitempty"` // Fields that failed validation
}

// InternalServerError represents errors that are unexpected and likely indicate a system issue.
type InternalServerError struct {
	Status  int    `json:"status"`          // HTTP status code
	Message string `json:"message"`         // Error message
	Stack   string `json:"stack,omitempty"` // Optional stack trace
}

// NotFoundError represents errors where a requested resource could not be found.
type NotFoundError struct {
	Status   int    `json:"status"`   // HTTP status code
	Message  string `json:"message"`  // Error message
	Resource string `json:"resource"` // The resource that could not be found
}

// APIError is a general-purpose error type for API responses.
type APIError struct {
	Status  int    `json:"status"`  // HTTP status code
	Message string `json:"message"` // Error message
}
