package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/internal/models"
	"github.com/pageza/chat-app/pkg/database"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var TokenExpiration string

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

	TokenExpiration = viper.GetString("TOKEN_EXPIRATION")

	fmt.Println("Debug TokenExpiration in auth test:", TokenExpiration)

	// Check if TokenExpiration is empty
	if TokenExpiration == "" {
		log.Fatalf("TokenExpiration is empty")
	}

	_, err := time.ParseDuration(TokenExpiration)
	if err != nil {
		log.Fatalf("Invalid token expiration duration in auth test: %v", err)
	}
	// Run the tests
	code := m.Run()

	// Teardown if needed

	os.Exit(code)
}

// Initialize or mock your database here
var db *gorm.DB

// CleanupTestUser removes the test user from the database.
func CleanupTestUser() {
	db.Where("username = ?", "testuser").Delete(&models.User{})
}

func initializeOrMockDatabase(t *testing.T) {
	// Get the DSN from the environment variable
	dsn := os.Getenv("TEST_DSN")
	if dsn == "" {
		t.Fatal("TEST_DSN environment variable not set")
	}

	// Connect to the PostgreSQL test database
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	database.DB = db

	// TODO: Run any database setup, schema migration, or seeding here
	if err := db.AutoMigrate(&models.User{}); err != nil {
		t.Fatalf("Could not migrate database: %v", err)
	}
}

// Initialize or mock your Redis client here
func initializeOrMockRedisClient() *redis.Client {
	// Initialize Redis client for testing
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // replace with your Redis server address and port
	})
	return rdb
}

// TestRegisterHandler tests the RegisterHandler function.
func TestRegisterHandler(t *testing.T) {
	testDSN := os.Getenv("TEST_DSN")
	if testDSN == "" {
		t.Fatalf("TEST_DSN environment variable not set")
	}
	fmt.Println("TEST_DSN is set to:", testDSN)
	// Initialize or mock database
	initializeOrMockDatabase(t)

	// Cleanup test user before running the test
	CleanupTestUser()

	user := models.User{
		Username: "testuser",
		Email:    "test@email.com",
		Password: "Test@123",
	}

	payload, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(RegisterHandler)

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Expected response code to be 201")

	// Check if user is created in the database
	var createdUser models.User
	db.Where("username = ?", user.Username).First(&createdUser)
	assert.Equal(t, user.Username, createdUser.Username, "User should be created in the database")
}

// TestLoginHandler tests the LoginHandler function.
func TestLoginHandler(t *testing.T) {
	// Initialize or mock database
	initializeOrMockDatabase(t)

	// Create a test user for login
	testUser := models.User{
		Username: "testuser",
		Email:    "test@email.com",
		Password: "Test@123",
	}
	db.Create(&testUser)

	user := models.User{
		Username: "testuser",
		Password: "Test@123",
	}

	payload, _ := json.Marshal(user)
	req, err := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(LoginHandler)

	handler.ServeHTTP(rr, req)
	resp, err := http.Get("your API endpoint")
	if err != nil {
		t.Fatalf("HTTP request failed: %v", err)
	}
	fmt.Println("Actual response:", resp)
	assert.Equal(t, http.StatusOK, rr.Code, "Expected response code to be 200")

	// Check if token is returned in the response
	assert.Contains(t, rr.Body.String(), "token", "Response should contain a token")
}

// TestLogoutHandler tests the LogoutHandler function.
func TestLogoutHandler(t *testing.T) {
	// Initialize or mock Redis client
	rdb := initializeOrMockRedisClient()

	req, err := http.NewRequest("POST", "/logout", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LogoutHandler(w, r, rdb)
	})

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Expected response code to be 200")

	// Check if the response indicates successful logout
	assert.Contains(t, rr.Body.String(), "Logged out successfully", "Response should indicate successful logout")
}
