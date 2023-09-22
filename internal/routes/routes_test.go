package routes_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pageza/chat-app/internal/auth"
	"github.com/pageza/chat-app/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/health", utils.HealthCheckHandler).Methods("GET")

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

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
}

// Additional test cases can be added here for different routes, middleware, etc.
