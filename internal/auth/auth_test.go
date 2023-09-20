package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pageza/chat-app/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestRegisterHandler tests the RegisterHandler function.
func TestRegisterHandler(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Email:    "test@email.com",
		Password: "Test@123",
	}

	payload, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected response code to be 201")
}

// TestLoginHandler tests the LoginHandler function.
func TestLoginHandler(t *testing.T) {
	user := models.User{
		Username: "testuser",
		Password: "Test@123",
	}

	payload, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected response code to be 200")
}

// TestLogoutHandler tests the LogoutHandler function.
func TestLogoutHandler(t *testing.T) {
	// Assuming rdb is your Redis client
	// rdb := redis.GetRedisClient()

	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LogoutHandler(w, r, nil) // Replace nil with your Redis client
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected response code to be 200")
}
