// Package jwt provides utility functions for generating, parsing, and managing JSON Web Tokens (JWT).
// This package is essential for user authentication and authorization in the application.

package jwt

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/errors"
	"github.com/pageza/chat-app/internal/models"
	"github.com/sirupsen/logrus"
)

// GenerateToken generates a JWT for a given user.
// The token will contain claims like the username and expiration time.
//
// Parameters:
// - user: The user for whom the token is generated
//
// Returns:
// - A signed JWT string
// - An error if something goes wrong
func GenerateToken(user models.User) (string, error) {
	// Parse the token expiration duration from the configuration
	expirationDuration, err := time.ParseDuration(config.TokenExpiration)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"duration": config.TokenExpiration,
		}).Fatalf("Invalid token expiration duration: %v", err)
	}

	// Calculate the expiration time for the token
	expirationTime := time.Now().Add(expirationDuration).Unix()

	// Create the claims for the token
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime,
		Issuer:    config.JwtIssuer,
		Subject:   user.Username,
	}

	// Generate the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JwtSecret))
}

// SetTokenCookie sets a JWT as an HttpOnly cookie.
//
// Parameters:
// - w: The http.ResponseWriter to write the cookie to
// - token: The JWT string to set as a cookie
func SetTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
}

// ParseToken parses a JWT string and returns the token object.
//
// Parameters:
// - tokenString: The JWT string to parse
//
// Returns:
// - A pointer to the parsed jwt.Token
// - An error if something goes wrong
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.WithFields(logrus.Fields{
				"alg": token.Header["alg"],
			}).Error("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtSecret), nil
	})
}

// GenerateTokenAndSetCookie generates a JWT for a user and sets it as a cookie.
//
// Parameters:
// - w: The http.ResponseWriter to write the cookie to
// - user: The user for whom the token is generated
func GenerateTokenAndSetCookie(w http.ResponseWriter, user models.User) {
	tokenString, err := GenerateToken(user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user": user.Username,
		}).Error("Could not generate token")
		apiErr := errors.NewAPIError(http.StatusInternalServerError, "Could not log in")
		errors.RespondWithError(w, apiErr)
		return
	}
	SetTokenCookie(w, tokenString)
}

// ClearTokenCookie clears the JWT cookie.
//
// Parameters:
// - w: The http.ResponseWriter to clear the cookie from
func ClearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
