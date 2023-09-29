// Package user contains functionalities related to user operations.
// TODO: Consider adding more user-related functionalities like updating profiles, password change, etc.

package user

import (
	"encoding/json"
	"net/http"

	"github.com/pageza/chat-app/internal/middleware"
	"github.com/pageza/chat-app/internal/models"
	"github.com/pageza/chat-app/pkg/database"
)

type UserDB interface {
	GetUserByID(userID string) (*models.User, error)
}

// TokenValidator is an interface for JWT token validation.
type TokenValidator interface {
	ValidateToken(r *http.Request) bool
}

// RealAuth is a concrete implementation of the TokenValidator interface.
type RealAuth struct{}

// ValidateToken calls the actual ValidateToken function from the middleware.
func (ra *RealAuth) ValidateToken(r *http.Request) bool {
	return middleware.ValidateToken(r)
}

// UserHandler contains dependencies for handling user-related requests.
type UserHandler struct {
	DB             database.Database
	TokenValidator TokenValidator
}

// UserInfoHandler handles the request to get user information.
// It expects the request to have valid JWT tokens in the header, which are validated by middleware.
func (uh *UserHandler) UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve username and email from request headers, which were set by ValidateMiddleware
	username := r.Header.Get("username")
	email := r.Header.Get("email")

	// If either username or email is empty, it means the JWT was not validated
	if username == "" || email == "" {
		// TODO: Consider logging the error for debugging
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create a user info JSON response
	userInfo := map[string]string{
		"username": username,
		"email":    email,
	}

	// Convert the map to JSON
	jsonResponse, err := json.Marshal(userInfo)
	if err != nil {
		// TODO: Consider logging the error for debugging
		http.Error(w, "Could not create user info response", http.StatusInternalServerError)
		return
	}

	// Set Content-Type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
