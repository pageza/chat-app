package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedis is a mock implementation of the Redis client
type MockRedis struct {
	mock.Mock
}

func (m *MockRedis) CheckRateLimit(ip string) (bool, error) {
	args := m.Called(ip)
	return args.Bool(0), args.Error(1)
}

func TestRateLimitMiddleware(t *testing.T) {
	// Create a request
	req, _ := http.NewRequest("GET", "/testRateLimit", nil)

	// Create a ResponseRecorder
	rr := httptest.NewRecorder()

	// Create a sample handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Create a mock Redis client
	mockRedis := new(MockRedis)

	// Test case 1: IP is within rate limit
	mockRedis.On("CheckRateLimit", mock.Anything).Return(true, nil)
	http.HandlerFunc(RateLimitMiddleware(handler).ServeHTTP).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	// Test case 2: IP has exceeded rate limit
	rr = httptest.NewRecorder() // Reset the ResponseRecorder
	mockRedis.On("CheckRateLimit", mock.Anything).Return(false, nil)
	http.HandlerFunc(RateLimitMiddleware(handler).ServeHTTP).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusTooManyRequests, rr.Code)

	// Test case 3: Internal error while checking rate limit
	rr = httptest.NewRecorder() // Reset the ResponseRecorder
	mockRedis.On("CheckRateLimit", mock.Anything).Return(false, errors.New("Internal error"))
	http.HandlerFunc(RateLimitMiddleware(handler).ServeHTTP).ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
