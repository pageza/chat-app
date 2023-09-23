// common_test.go

package common

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
	status := http.StatusBadRequest
	message := "Bad Request"

	apiErr := NewAPIError(status, message)

	assert.Equal(t, status, apiErr.Status, "Status should match")
	assert.Equal(t, message, apiErr.Message, "Message should match")
}

// TestRespondWithError tests the RespondWithError function
func TestRespondWithError(t *testing.T) {
	status := http.StatusNotFound
	message := "Not Found"

	apiErr := NewAPIError(status, message)

	_, err := http.NewRequest("GET", "/nonexistent", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	RespondWithError(rr, apiErr)

	assert.Equal(t, status, rr.Code, "Status code should match")

	var response APIError
	err = json.NewDecoder(rr.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, status, response.Status, "Response status should match")
	assert.Equal(t, message, response.Message, "Response message should match")
}
