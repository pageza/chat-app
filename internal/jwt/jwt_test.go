package jwt

import (
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/models"
	"github.com/stretchr/testify/assert"
)

func init() {
	config.JwtSecret = "testSecret"
	config.TokenExpiration = "1h"
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

	cookie := w.Header().Get("Set-Cookie")
	assert.Contains(t, cookie, "token=")
	assert.Contains(t, cookie, "Max-Age=-1")
}
