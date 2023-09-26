package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pageza/chat-app/internal/auth"
	"github.com/pageza/chat-app/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestRegisterHandler tests the RegisterHandler function
func TestRegisterHandler(t *testing.T) {
	// Initialize your dependencies here, like database, cache, etc.
	// For now, let's assume you have a function called `setupDependencies()`
	// that does this for you.
	// setupDependencies()

	// Test successful registration
	t.Run("successful registration", func(t *testing.T) {
		// Create a new user payload
		user := models.User{
			Username: "testuser",
			Email:    "test@email.com",
			Password: "Test@123",
		}
		payload, _ := json.Marshal(user)

		// Create a new request
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()

		// Call the RegisterHandler
		handler := http.HandlerFunc(auth.RegisterHandler)
		handler.ServeHTTP(rr, req)

		// Check the status code and the response
		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "User registered successfully")
	})

	// Test registration with existing username
	t.Run("existing username", func(t *testing.T) {
		// Create a new user payload with an existing username
		user := models.User{
			Username: "existinguser",
			Email:    "new@email.com",
			Password: "Test@123",
		}
		payload, _ := json.Marshal(user)

		// Create a new request
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()

		// Call the RegisterHandler
		handler := http.HandlerFunc(auth.RegisterHandler)
		handler.ServeHTTP(rr, req)

		// Check the status code and the response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Username already exists")
	})

	// Test registration with existing email
	t.Run("existing email", func(t *testing.T) {
		// Create a new user payload with an existing email
		user := models.User{
			Username: "newuser",
			Email:    "existing@email.com",
			Password: "Test@123",
		}
		payload, _ := json.Marshal(user)

		// Create a new request
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()

		// Call the RegisterHandler
		handler := http.HandlerFunc(auth.RegisterHandler)
		handler.ServeHTTP(rr, req)

		// Check the status code and the response
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "Email already exists")
	})
}
