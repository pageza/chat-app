package routes_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/internal/auth"
	"github.com/pageza/chat-app/internal/utils"
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

func TestHealthCheckHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/health", utils.HealthCheckHandler).Methods("GET")

	req, err := http.NewRequest("GET", "/health", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "OK", rr.Body.String())
}

func TestRegisterHandler(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/register", auth.RegisterHandler).Methods("POST")

	payload := `{"username": "test", "password": "password", "email": "test@example.com"}`
	req, err := http.NewRequest("POST", "/register", strings.NewReader(payload))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	// Add more assertions based on what your RegisterHandler is supposed to return
}

// Additional test cases can be added here for different routes, middleware, etc.
