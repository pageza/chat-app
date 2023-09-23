package jwt

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/models"
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

func init() {
	config.JwtSecret = "testSecret"
	// config.TokenExpiration = "1h"
	config.JwtIssuer = "testIssuer"
}

func TestGenerateToken(t *testing.T) {
	user := models.User{Username: "testUser"}
	token, err := GenerateToken(user)
	assert.Nil(t, err)

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecret), nil
	})

	assert.Nil(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestSetTokenCookie(t *testing.T) {
	w := httptest.NewRecorder()
	SetTokenCookie(w, "testToken")

	cookie := w.Header().Get("Set-Cookie")
	assert.Contains(t, cookie, "token=testToken")
}

func TestParseToken(t *testing.T) {
	user := models.User{Username: "testUser"}
	token, _ := GenerateToken(user)

	parsedToken, err := ParseToken(token)
	assert.Nil(t, err)
	assert.True(t, parsedToken.Valid)
}

func TestGenerateTokenAndSetCookie(t *testing.T) {
	user := models.User{Username: "testUser"}
	w := httptest.NewRecorder()

	GenerateTokenAndSetCookie(w, user)

	cookie := w.Header().Get("Set-Cookie")
	assert.Contains(t, cookie, "token=")
}

func TestClearTokenCookie(t *testing.T) {
	w := httptest.NewRecorder()
	ClearTokenCookie(w)
	cookie := w.Result().Cookies()[0] // Assuming it's the first cookie
	if cookie.MaxAge != -1 {
		t.Errorf("Expected MaxAge to be -1, got %d", cookie.MaxAge)
	}
	if cookie.HttpOnly != true {
		t.Errorf("Expected HttpOnly to be true")
	}
	fmt.Println("Cookie string:", cookie.String())

	// Convert the cookie to its string representation
	cookieStr := cookie.String()
	assert.Contains(t, cookieStr, "token=")
	assert.Contains(t, cookieStr, "Max-Age=0")
}
