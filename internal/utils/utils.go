// Package utils contains utility functions that can be used across the application.

package utils

import (
	"fmt"
	"net/http"

	"github.com/pageza/chat-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// HealthCheckHandler is an HTTP handler function that checks the health of the server.
// It responds with an HTTP 200 status code and a message indicating that the server is up and running.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// Set the status code to 200 OK
	w.WriteHeader(http.StatusOK)

	// Write the health check message to the response
	fmt.Fprintf(w, "Server is up and running!")
}

// ValidateUser compares the hashed password stored in the User model with the provided plaintext password.
// It uses bcrypt's CompareHashAndPassword function for secure password comparison.
// Returns an error if the passwords do not match.
func ValidateUser(user *models.User, password string) error {
	// Compare the hashed password in the user model with the provided plaintext password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	// Return the result of the comparison (nil if passwords match, error otherwise)
	return err
}
