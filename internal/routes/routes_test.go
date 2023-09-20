package routes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pageza/chat-app/internal/auth"
	"github.com/pageza/chat-app/internal/middleware"
	"github.com/pageza/chat-app/internal/user"
	"github.com/pageza/chat-app/internal/utils"
	"github.com/stretchr/testify/assert"
)

// TestHealthCheckHandler tests the health check route
func TestHealthCheckHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/health", utils.HealthCheckHandler).Methods("GET")

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())

	// TODO: Add additional tests for different scenarios, if needed.
}

// TestRegisterHandler tests the register route
func TestRegisterHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/register", auth.RegisterHandler).Methods("POST")

	payload := `{"username": "test", "password": "password", "email": "test@example.com"}`
	req, err := http.NewRequest("POST", "/register", strings.NewReader(payload))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	// Add more assertions based on what your RegisterHandler is supposed to return

	// TODO: Add tests for invalid payload
	// TODO: Add tests for existing username or email
	// TODO: Add tests for server errors
}

// TestUserInfoHandler tests the user info route with authentication middleware
func TestUserInfoHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/userinfo", middleware.AuthMiddleware(user.UserInfoHandler)).Methods("GET")

	req, err := http.NewRequest("GET", "/userinfo", nil)
	// Add Authorization header to mock a valid JWT token
	req.Header.Set("Authorization", "Bearer mockValidToken")
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	// Add more assertions based on what your UserInfoHandler is supposed to return

	// TODO: Add tests for invalid or missing JWT token
	// TODO: Add tests for expired JWT token
	// TODO: Add tests for server errors
}

// TODO: Add tests for other routes like /chat, /send, /receive, /login, /logout, etc.
