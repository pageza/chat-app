package server_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/pageza/chat-app/internal/redis"
	"github.com/pageza/chat-app/internal/routes"
	"github.com/stretchr/testify/assert"
)

// TestStartServer tests if the server starts successfully.
func TestStartServer(t *testing.T) {
	// Set environment variable for server port
	os.Setenv("SERVER_PORT", "8080")

	// Initialize the Gorilla Mux router
	r := mux.NewRouter()

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Handle the request using the router
	r.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")
}

// TestMiddlewareApplied tests if the middleware is correctly applied to the routes.
func TestMiddlewareApplied(t *testing.T) {
	// Initialize the Gorilla Mux router
	r := mux.NewRouter()

	// Get the Redis client instance
	rdb := redis.GetRedisClient()

	// Add middleware and routes
	routes.InitializeRoutes(r, rdb)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/some-route", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Handle the request using the router
	r.ServeHTTP(rr, req)

	// Check if middleware logic is executed (e.g., check headers, status code, etc.)
	assert.Equal(t, http.StatusUnauthorized, rr.Code, "Status code should be Unauthorized")
}

// TestCORS tests if CORS settings are correctly applied.
func TestCORS(t *testing.T) {
	// Initialize the Gorilla Mux router
	r := mux.NewRouter()

	// Get the Redis client instance
	rdb := redis.GetRedisClient()

	// Add middleware and routes
	routes.InitializeRoutes(r, rdb)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/some-route", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Handle the request using the router
	r.ServeHTTP(rr, req)

	// Check if CORS headers are present
	assert.Equal(t, "*", rr.Header().Get("Access-Control-Allow-Origin"))
}

// TestRoutes tests if the routes are correctly initialized.
func TestRoutes(t *testing.T) {
	r := mux.NewRouter()
	// Get the Redis client instance
	rdb := redis.GetRedisClient()

	// Add middleware and routes
	routes.InitializeRoutes(r, rdb)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Handle the request using the router
	r.ServeHTTP(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be OK")
}
