// Package user contains functionalities related to user operations.
// TODO: Consider adding more user-related functionalities like updating profiles, password change, etc.

package user

import (
	"encoding/json"
	"net/http"
)

// UserInfoHandler handles the request to get user information.
// It expects the request to have valid JWT tokens in the header, which are validated by middleware.
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve username and email from request headers, which were set by ValidateMiddleware
	username := r.Header.Get("username")
	email := r.Header.Get("email")

	// If either username or email is empty, it means the JWT was not validated
	if username == "" || email == "" {
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
		http.Error(w, "Could not create user info response", http.StatusInternalServerError)
		return
	}

	// Set Content-Type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
