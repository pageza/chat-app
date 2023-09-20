package config

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitializeCORS(t *testing.T) {
	// Initialize CORS middleware
	corsMiddleware := InitializeCORS()

	// Create a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		corsMiddleware.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
			// This is your actual handler
			w.Write([]byte("Hello, world!"))
		})
	}))
	defer ts.Close()

	// Make a request to the test server
	client := &http.Client{}
	req, err := http.NewRequest("GET", ts.URL, nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// Add Origin header to simulate a CORS request
	req.Header.Set("Origin", "http://example.com")

	// Perform the request
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check if CORS headers are set
	assert.Equal(t, "http://example.com", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
	assert.Equal(t, "Authorization", resp.Header.Get("Access-Control-Expose-Headers"))
}
