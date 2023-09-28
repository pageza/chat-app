package auth_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/pageza/chat-app/internal/auth"
	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/jwt"
	"github.com/pageza/chat-app/internal/models"
	"github.com/pageza/chat-app/internal/utils"
)

type MockJwt struct {
	mock.Mock
}

func (m *MockJwt) GenerateToken(user models.User) (string, string, error) {
	args := m.Called(user)
	return args.String(0), args.String(1), args.Error(2)
}

type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) InitializeDB() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

func (m *MockDatabase) AutoMigrateDB() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDatabase) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDatabase) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	user, ok := args.Get(0).(*models.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockDatabase) UpdateLastLoginTime(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDatabase) HandleFailedLoginAttempt(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDatabase) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
}

type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	return args.Get(0).(*redis.StatusCmd)
}

func TestHandlers(t *testing.T) {
	mockDB := new(MockDatabase)
	mockRedisClient := new(MockRedisClient)

	mockDB.On("InitializeDB").Return(new(gorm.DB), nil)
	mockDB.On("AutoMigrateDB").Return(nil)
	mockDB.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)
	mockDB.On("GetUserByUsername", mock.AnythingOfType("string")).Return(new(models.User), nil)
	mockDB.On("UpdateLastLoginTime", mock.AnythingOfType("*models.User")).Return(nil)
	mockDB.On("HandleFailedLoginAttempt", mock.AnythingOfType("*models.User")).Return(nil)
	mockDB.On("Where", mock.Anything, mock.Anything).Return(new(gorm.DB))

	mockRedisClient.On("Set", mock.Anything, mock.Anything, mock.Anything, mock.AnythingOfType("time.Duration")).Return(new(redis.StatusCmd))

	authHandler := &auth.AuthHandler{
		DB: mockDB,
	}

	t.Run("Test LogoutHandler", func(t *testing.T) {
		user := models.User{
			Username: "testuser",
		}

		accessToken, _, err := jwt.GenerateToken(user)
		if err != nil {
			t.Fatalf("Could not generate token: %v", err)
		}

		req, _ := http.NewRequest("POST", "/logout", nil)
		rr := httptest.NewRecorder()

		req.Header.Set("Authorization", "Bearer "+accessToken)

		authHandler.LogoutHandler(rr, req, mockRedisClient)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "Logged out successfully")
	})

	t.Run("Test RegisterHandler", func(t *testing.T) {
		// Create a user object
		user := models.User{
			Username: "newuser",
			Password: "newpassword",
		}

		// Convert user object to JSON
		payload, _ := json.Marshal(user)

		// Create a new request
		req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()

		// Call the RegisterHandler
		authHandler.RegisterHandler(rr, req)

		// Check the status code and the response
		assert.Equal(t, http.StatusCreated, rr.Code)
		// Add more assertions based on your requirements
	})
}

func TestLoginHandler(t *testing.T) {
	config.Initialize()
	// Initialize dependencies and mock objects
	dbMock := new(MockDatabase)
	jwtMock := new(MockJwt)
	authHandler := &auth.AuthHandler{
		DB:           dbMock,
		JwtGenerator: jwtMock,
	}

	// Define specific errors for use in tests
	userNotFoundError := errors.New("user not found")

	t.Run("Valid credentials", func(t *testing.T) {
		fmt.Println("Entering test: Valid credentials")
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
		user := models.User{
			ID:       1,
			Username: "testuser",
			Password: string(hashedPassword),
		}

		fmt.Println("Debug: Hashed password:", string(hashedPassword))

		err := utils.ValidateUser(&user, "testpassword")
		if err != nil {
			t.Fatalf("Password validation failed: %v", err)
		}
		fmt.Println("Running test: Valid credentials- ", user)

		dbMock.On("GetUserByUsername", "testuser").Return(&user, nil)
		// Add this debugging log
		fmt.Println("Debug: User returned from mock DB:", user)
		jwtMock.On("GenerateToken", user).Return("accessToken", "refreshToken", nil)
		jwtMock.On("GenerateToken", user).Return("accessToken", "refreshToken", nil)
		// Create a payload with the plaintext password for the HTTP request
		requestPayload := models.User{
			Username: "testuser",
			Password: "testpassword", // Plaintext password
		}
		payload, _ := json.Marshal(requestPayload)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()
		authHandler.LoginHandler(rr, req)
		assert.Equal(t, http.StatusOK, rr.Code)
	})

	// TODO: Revisit this test for token generation failure when refining the authentication module.
	// t.Run("Token generation failure", func(t *testing.T) {
	// 	// Print debug information to indicate the start of this test case
	// 	fmt.Println("Entering test: Token generation failure")
	// 	fmt.Println("Running test: Token generation failure")

	// 	// Create a mock user object
	// 	user := models.User{
	// 		Username: "testuser",
	// 		Password: "testpassword",
	// 	}

	// 	// Generate a hashed password that matches the plaintext password
	// 	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("testpassword"), bcrypt.DefaultCost)
	// 	dbUser := &models.User{
	// 		Username: "testuser",
	// 		Password: string(hashedPassword),
	// 	}

	// 	// Mock the database call to return the mock user object when GetUserByUsername is called
	// 	dbMock.On("GetUserByUsername", "testuser").Return(dbUser, nil) // <-- Change here

	// 	// Mock the JWT token generation to return an error
	// 	tokenGenerationError := errors.New("Token generation failed")
	// 	jwtMock.On("GenerateToken", *dbUser).Return("", "", tokenGenerationError) // <-- Change here

	// 	// Create a payload for the HTTP request
	// 	payload, _ := json.Marshal(user)

	// 	// Create a new HTTP request for the login endpoint
	// 	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))

	// 	// Create a ResponseRecorder to record the HTTP response
	// 	rr := httptest.NewRecorder()

	// 	// Call the LoginHandler function
	// 	authHandler.LoginHandler(rr, req)

	// 	// Assert that the mock was called
	// 	assert.True(t, jwtMock.AssertExpectations(t), "GenerateToken mock not called")

	// 	// Assert that the HTTP status code should be 500 (Internal Server Error)
	// 	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	// })

	t.Run("Invalid credentials", func(t *testing.T) {
		fmt.Println("Entering test: Invalid credentials")
		fmt.Println("Running test: Invalid credentials")
		user := models.User{
			Username: "wronguser",
			Password: "wrongpassword",
		}
		dbMock.On("GetUserByUsername", "wronguser").Return(nil, userNotFoundError)

		payload, _ := json.Marshal(user)
		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(payload))
		rr := httptest.NewRecorder()
		authHandler.LoginHandler(rr, req)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})
}
