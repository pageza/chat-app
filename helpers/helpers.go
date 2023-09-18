package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/middleware"
	"github.com/pageza/chat-app/models"
	"golang.org/x/crypto/bcrypt"
)

func ValidateUser(user *models.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err
}
func SetTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
}
func HandleError(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})
}
func BlacklistToken(rdb *redis.Client, tokenString string, expirationTime int64) error {
	return rdb.Set(context.TODO(), tokenString, "blacklisted", time.Until(time.Unix(expirationTime, 0))).Err()
}

// SendJSONResponse sends a JSON response to the client
func SendJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func GenerateTokenAndSetCookie(w http.ResponseWriter, user models.User) {
	tokenString, err := middleware.GenerateToken(user)
	if err != nil {
		HandleError(w, "Could not log in", http.StatusInternalServerError)
		return
	}
	SetTokenCookie(w, tokenString)
}

// ClearTokenCookie clears the token cookie
func ClearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // MaxAge<0 means delete cookie now
	})
}
