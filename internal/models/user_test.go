package models

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.Println("Starting TestMain...") // This should appear in your test output

	// Load .env.test
	if err := godotenv.Load("/home/zach/projects/chat-app/.env"); err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

	// Initialize Viper
	viper.SetConfigFile("/home/zach/projects/chat-app/internal/config/config.yaml") // Adjust the path as needed
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config.yaml: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Teardown if needed

	os.Exit(code)
}

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
