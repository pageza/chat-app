package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pageza/chat-app/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestRecoveryMiddleware(t *testing.T) {
	// Create a request
	req, _ := http.NewRequest("GET", "/testRecovery", nil)

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a sample handler that panics
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("This is a test panic")
	})

	// Wrap the handler with RecoveryMiddleware
	http.HandlerFunc(RecoveryMiddleware(http.HandlerFunc(handler)).ServeHTTP).ServeHTTP(rr, req)

	// Check the status code and body
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var apiErr errors.APIError
	_ = json.Unmarshal(rr.Body.Bytes(), &apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
	assert.Equal(t, "Internal Server Error", apiErr.Message)
}
