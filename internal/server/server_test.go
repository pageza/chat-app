package server_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/internal/redis"
	"github.com/pageza/chat-app/internal/routes"
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

// Additional test cases can be added here for different routes, middleware, etc.
