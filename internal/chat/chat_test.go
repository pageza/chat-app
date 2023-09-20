// chat_test.go

package chat

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
