package utils

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/internal/models"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
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

// TestHealthCheckHandler tests the HealthCheckHandler function.
func TestHealthCheckHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HealthCheckHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "Server is up and running!"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

// TestValidateUser tests the ValidateUser function.
func TestValidateUser(t *testing.T) {
	password := "securePassword123!"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Password: string(hashedPassword),
	}

	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"Valid password", password, false},
		{"Invalid password", "wrongPassword", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUser(user, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
