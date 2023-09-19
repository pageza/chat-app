package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/internal/common"
	"github.com/pageza/chat-app/internal/errors" // <-- Updated import
	"github.com/pageza/chat-app/internal/models"
	"github.com/sirupsen/logrus"
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

func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(common.JwtSecret), nil
	})
}

func BlacklistToken(rdb *redis.Client, tokenString string, expirationTime int64) error {
	const maxRetries = 3
	var currentRetry = 0

	for currentRetry < maxRetries {
		err := rdb.Set(context.TODO(), tokenString, "blacklisted", time.Until(time.Unix(expirationTime, 0))).Err()
		if err == nil {
			return nil // Operation was successful, return
		}

		currentRetry++
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	if currentRetry == maxRetries {
		logrus.Printf("Max retries reached, could not blacklist the token")
		return fmt.Errorf("max retries reached, could not blacklist the token")
	}

	return nil // This line is technically unreachable but added for completeness
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}

func GenerateTokenAndSetCookie(w http.ResponseWriter, user models.User) {
	tokenString, err := common.GenerateToken(user)
	if err != nil {
		apiErr := errors.NewAPIError(http.StatusInternalServerError, "Could not log in") // <-- Updated line
		errors.RespondWithError(w, apiErr)                                               // <-- Updated line
		return
	}
	SetTokenCookie(w, tokenString)
}

func ClearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
