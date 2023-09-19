package common

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pageza/chat-app/internal/models"
	"github.com/sirupsen/logrus"
)

var (
	JwtSecret       string
	JwtIssuer       string
	RedisAddr       string
	TokenExpiration string
)

func Initialize() {
	JwtSecret = os.Getenv("JWT_SECRET")
	JwtIssuer = os.Getenv("JWT_ISSUER")
	RedisAddr = os.Getenv("REDIS_ADDR")
	TokenExpiration = os.Getenv("TOKEN_EXPIRATION")

	if JwtSecret == "" || JwtIssuer == "" || RedisAddr == "" || TokenExpiration == "" {
		logrus.Fatal("Environment variables are not set")
	}

	_, err := time.ParseDuration(TokenExpiration)
	if err != nil {
		logrus.Fatalf("Invalid token expiration duration: %v", err)
	}
}

func GenerateToken(user models.User) (string, error) {
	expirationDuration, err := time.ParseDuration(TokenExpiration)
	if err != nil {
		logrus.Fatalf("Invalid token expiration duration: %v", err)
	}

	expirationTime := time.Now().Add(expirationDuration).Unix()

	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime,
		Issuer:    JwtIssuer,
		Subject:   user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(JwtSecret))
}

// APIError represents an error that can be sent in an API response
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

// NewAPIError creates a new APIError
func NewAPIError(status int, message string) *APIError {
	return &APIError{
		Status:  status,
		Message: message,
	}
}

// RespondWithError sends an API response with a error message and status code
func RespondWithError(w http.ResponseWriter, err *APIError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}
