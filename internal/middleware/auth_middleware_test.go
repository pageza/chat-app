package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/pageza/chat-app/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestAuthMiddleware(t *testing.T) {
	// Generate a sample token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: 15000,
		Issuer:    "test",
	})
	tokenString, _ := token.SignedString([]byte(config.JwtSecret))

	// Create a request with the token
	req, _ := http.NewRequest("GET", "/testAuth", nil)
	req.AddCookie(&http.Cookie{
		Name:  "token",
		Value: tokenString,
	})

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a sample handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Wrap the handler with AuthMiddleware
	http.HandlerFunc(AuthMiddleware(handler)).ServeHTTP(rr, req)

	// Check the status code and body
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

// Additional test cases can be added here for invalid tokens, expired tokens, etc.
