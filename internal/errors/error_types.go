// error_types.go

package errors

// AuthenticationError for errors related to user authentication, such as invalid credentials or expired tokens.
type AuthenticationError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Reason  string `json:"reason,omitempty"`
}

// RateLimitError for errors related to exceeding rate limits.
type RateLimitError struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
}

// DatabaseError for errors related to database operations, such as failed queries or connection issues.
type DatabaseError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Query   string `json:"query,omitempty"`
}

// ValidationError for errors related to invalid input or failed validation.
type ValidationError struct {
	Status  int      `json:"status"`
	Message string   `json:"message"`
	Fields  []string `json:"fields,omitempty"`
}

// InternalServerError for errors that are unexpected and likely indicate a bug or system issue.
type InternalServerError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Stack   string `json:"stack,omitempty"`
}

// NotFoundError for errors where a resource could not be found.
type NotFoundError struct {
	Status   int    `json:"status"`
	Message  string `json:"message"`
	Resource string `json:"resource"`
}

// APIError represents an error that can be sent in an API response
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
