// chat_test.go

package chat

import (
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

// TestChatHandler tests the ChatHandler function
func TestChatHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/chat", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ChatHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")
	assert.Equal(t, "Chat endpoint", rr.Body.String(), "Response body should match")
}

// TestSendMessageHandler tests the SendMessageHandler function
func TestSendMessageHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/send", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(SendMessageHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")
	assert.Equal(t, "Send message endpoint", rr.Body.String(), "Response body should match")
}

// TestReceiveMessageHandler tests the ReceiveMessageHandler function
func TestReceiveMessageHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/receive", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ReceiveMessageHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")
	assert.Equal(t, "Receive message endpoint", rr.Body.String(), "Response body should match")
}
