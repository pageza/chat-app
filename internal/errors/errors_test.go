package errors

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
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

// TestNewAPIError tests the NewAPIError function
func TestNewAPIError(t *testing.T) {
	err := NewAPIError(http.StatusNotFound, "Resource not found")
	assert.Equal(t, http.StatusNotFound, err.Status)
	assert.Equal(t, "Resource not found", err.Message)
}

// TestRespondWithError tests the RespondWithError function
func TestRespondWithError(t *testing.T) {
	err := NewAPIError(http.StatusNotFound, "Resource not found")
	w := httptest.NewRecorder()
	RespondWithError(w, err)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var apiErr APIError
	json.Unmarshal(w.Body.Bytes(), &apiErr)
	assert.Equal(t, err.Status, apiErr.Status)
	assert.Equal(t, err.Message, apiErr.Message)
}

// TestRespondWithCustomError tests the RespondWithCustomError function
func TestRespondWithCustomError(t *testing.T) {
	customErr := &AuthenticationError{
		Status:  http.StatusUnauthorized,
		Message: "Invalid credentials",
		Reason:  "Username or password is incorrect",
	}
	w := httptest.NewRecorder()
	RespondWithCustomError(w, customErr)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))

	var authErr AuthenticationError
	json.Unmarshal(w.Body.Bytes(), &authErr)
	assert.Equal(t, customErr.Status, authErr.Status)
	assert.Equal(t, customErr.Message, authErr.Message)
	assert.Equal(t, customErr.Reason, authErr.Reason)
}
