package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/internal/errors"
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

func TestRecoveryMiddleware(t *testing.T) {
	// Create a request
	req, _ := http.NewRequest("GET", "/testRecovery", nil)

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a sample handler that panics
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("This is a test panic")
	})

	// Wrap the handler with RecoveryMiddleware
	http.HandlerFunc(RecoveryMiddleware(http.HandlerFunc(handler)).ServeHTTP).ServeHTTP(rr, req)

	// Check the status code and body
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var apiErr errors.APIError
	_ = json.Unmarshal(rr.Body.Bytes(), &apiErr)
	assert.Equal(t, http.StatusInternalServerError, apiErr.Status)
	assert.Equal(t, "Internal Server Error", apiErr.Message)
}
