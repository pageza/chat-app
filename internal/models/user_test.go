package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserValidation(t *testing.T) {
	// Test cases
	tests := []struct {
		username string
		email    string
		password string
		isValid  bool
	}{
		{"john", "john@example.com", "Password@123", true},
		{"jo", "john@example.com", "Password@123", false},  // Invalid username
		{"john", "johnexample.com", "Password@123", false}, // Invalid email
		{"john", "john@example.com", "password", false},    // Invalid password
	}

	for _, test := range tests {
		// Create a new User instance
		user := &User{
			Username: test.username,
			Email:    test.email,
			Password: test.password,
		}

		// Validate the user
		err := user.Validate()

		// Check if the validation result matches the expected outcome
		if test.isValid {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}
