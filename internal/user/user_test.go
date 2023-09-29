package user_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pageza/chat-app/internal/models"
	"github.com/pageza/chat-app/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockAuth struct {
	mock.Mock
}

func (a *MockAuth) ValidateToken(r *http.Request) bool {
	args := a.Called(r)
	return args.Bool(0)
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) InitializeDB() (*gorm.DB, error) {
	args := m.Called()
	return args.Get(0).(*gorm.DB), args.Error(1)
}

func (m *MockDB) AutoMigrateDB() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockDB) CreateUser(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) GetUserByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	user, ok := args.Get(0).(*models.User)
	if !ok {
		return nil, args.Error(1)
	}
	return user, args.Error(1)
}

func (m *MockDB) UpdateLastLoginTime(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) HandleFailedLoginAttempt(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	callArgs := m.Called(query, args)
	return callArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) GetUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := m.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func TestUserInfoHandler(t *testing.T) {
	// Initialize mocks
	dbMock := new(MockDB)
	authMock := new(MockAuth)

	// Initialize the handler
	userHandler := &user.UserHandler{
		DB: dbMock,
	}
	t.Run("Successful User Information Retrieval", func(t *testing.T) {
		// Mock the JWT validation to return a valid user ID
		authMock.On("ValidateToken", mock.Anything).Return("testuser", nil)

		// Mock the database call to return a mock user object
		mockUser := &models.User{
			Username: "testuser",
			Email:    "test@email.com",
		}
		dbMock.On("GetUserByID", "testuser").Return(mockUser, nil)

		// Create a new HTTP request for the user info endpoint
		req, _ := http.NewRequest("GET", "/user/info", nil)

		// Create a ResponseRecorder to record the HTTP response
		rr := httptest.NewRecorder()

		// Call the UserInfoHandler function
		userHandler.UserInfoHandler(rr, req)

		// Assert that the HTTP status code should be 200 (OK)
		assert.Equal(t, http.StatusOK, rr.Code)

		// Assert that the JSON response contains the correct user information
		var responseUser models.User
		json.Unmarshal(rr.Body.Bytes(), &responseUser)
		assert.Equal(t, mockUser, &responseUser)
	})

	t.Run("Unauthorized User", func(t *testing.T) {
		// Mock the JWT validation to return an error
		authMock.On("ValidateToken", mock.Anything).Return("", errors.New("Unauthorized"))

		// Create a new HTTP request for the user info endpoint
		req, _ := http.NewRequest("GET", "/user/info", nil)

		// Create a ResponseRecorder to record the HTTP response
		rr := httptest.NewRecorder()

		// Call the UserInfoHandler function
		userHandler.UserInfoHandler(rr, req)

		// Assert that the HTTP status code should be 401 (Unauthorized)
		assert.Equal(t, http.StatusUnauthorized, rr.Code)
	})

	t.Run("Internal Server Error", func(t *testing.T) {
		// Mock the JWT validation to return a valid user ID
		authMock.On("ValidateToken", mock.Anything).Return("testuser", nil)

		// Inject an error in the JSON marshaling process by returning an invalid user object from the database
		dbMock.On("GetUserByID", "testuser").Return(nil, errors.New("Internal Server Error"))

		// Create a new HTTP request for the user info endpoint
		req, _ := http.NewRequest("GET", "/user/info", nil)

		// Create a ResponseRecorder to record the HTTP response
		rr := httptest.NewRecorder()

		// Call the UserInfoHandler function
		userHandler.UserInfoHandler(rr, req)

		// Assert that the HTTP status code should be 500 (Internal Server Error)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
