package user_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/internal/user"
	"github.com/spf13/viper"
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

func TestUserInfoHandler(t *testing.T) {
	tests := []struct {
		name           string
		username       string
		email          string
		expectedStatus int
		expectedBody   map[string]string
	}{
		{
			name:           "Valid User",
			username:       "john",
			email:          "john@example.com",
			expectedStatus: http.StatusOK,
			expectedBody:   map[string]string{"username": "john", "email": "john@example.com"},
		},
		{
			name:           "Missing Username",
			username:       "",
			email:          "john@example.com",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
		{
			name:           "Missing Email",
			username:       "john",
			email:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "/userinfo", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("username", tt.username)
			req.Header.Set("email", tt.email)

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(user.UserInfoHandler)

			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v", status, tt.expectedStatus)
			}

			if tt.expectedBody != nil {
				expectedResponse, _ := json.Marshal(tt.expectedBody)
				if rr.Body.String() != string(expectedResponse) {
					t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), string(expectedResponse))
				}
			}
		})
	}

}
