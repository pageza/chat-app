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

func GenerateToken(user models.User) (string, error) {
	expirationDuration, err := time.ParseDuration(config.TokenExpiration)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"duration": config.TokenExpiration,
		}).Fatalf("Invalid token expiration duration: %v", err)
	}

	expirationTime := time.Now().Add(expirationDuration).Unix()

	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime,
		Issuer:    config.JwtIssuer,
		Subject:   user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.JwtSecret))
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
			logrus.WithFields(logrus.Fields{
				"alg": token.Header["alg"],
			}).Error("Unexpected signing method")
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtSecret), nil
	})
}

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

func ClearTokenCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
}
