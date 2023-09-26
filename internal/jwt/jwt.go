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
func init() {
	config.Initialize()
}

func GenerateToken(user models.User) (string, string, error) {
	// Parse the token expiration duration from the configuration for access token
	fmt.Println("Debug TokenExpiration in GenerateToken:", config.TokenExpiration)
	logrus.Infof("JWT - TokenExpiration: %s", config.TokenExpiration)

	expirationDuration, err := time.ParseDuration(config.TokenExpiration)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"duration": config.TokenExpiration,
		}).Fatalf("Invalid token expiration duration jwt: %v", err)
		return "", "", err
	}

	// Calculate the expiration time for the access token
	expirationTime := time.Now().Add(expirationDuration).Unix()

	// Create the claims for the access token
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime,
		Issuer:    config.JwtIssuer,
		Subject:   user.Username,
	}

	// Generate the access token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return "", "", err
	}

	// Generate refresh token with longer expiration (e.g., 24 hours)
	refreshExpirationTime := time.Now().Add(24 * time.Hour).Unix()
	refreshClaims := &jwt.StandardClaims{
		ExpiresAt: refreshExpirationTime,
		Issuer:    config.JwtIssuer,
		Subject:   user.Username,
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(config.JwtSecret))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshTokenString, nil
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
	accessToken, refreshToken, err := GenerateToken(user) // Updated this line
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user": user.Username,
		}).Error("Could not generate token")
		apiErr := errors.NewAPIError(http.StatusInternalServerError, "Could not log in")
		errors.RespondWithError(w, apiErr)
		return
	}
	SetTokenCookie(w, accessToken) // Assuming this sets the access token cookie

	// Set the refresh token as a cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Expires:  time.Now().Add(48 * time.Hour), // Set your desired expiration time
		HttpOnly: true,
	})
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
