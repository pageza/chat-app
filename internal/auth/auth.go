package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/internal/errors" // <-- Updated import
	jwtI "github.com/pageza/chat-app/internal/jwt"
	"github.com/pageza/chat-app/internal/models"
	redisI "github.com/pageza/chat-app/internal/redis"
	"github.com/pageza/chat-app/internal/utils"
	"github.com/pageza/chat-app/pkg/database"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Warn("Invalid request payload")
		errors.RespondWithError(w, errors.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	const maxRetries = 3
	var currentRetry = 0
	var result *gorm.DB

	for currentRetry < maxRetries {
		result = database.DB.Create(&user)
		if result.Error == nil {
			break
		}

		currentRetry++
		time.Sleep(2 * time.Second)
	}

	if currentRetry == maxRetries || result.Error != nil {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Warn("Could not register user after max retries")
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Could not register user"))
		return
	}

	tokenString, err := jwtI.GenerateToken(user)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Could not log in"))
		return
	}

	jwtI.SetTokenCookie(w, tokenString)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User successfully registered and logged in")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Warn("Invalid request payload")
		errors.RespondWithError(w, errors.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	const maxRetries = 3
	var currentRetry = 0
	var dbUser models.User

	for currentRetry < maxRetries {
		database.DB.Where("username = ?", user.Username).First(&dbUser)
		if dbUser.ID != 0 {
			break
		}

		currentRetry++
		time.Sleep(2 * time.Second)
	}

	if currentRetry == maxRetries || dbUser.ID == 0 {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Warn("Invalid credentials after max retries")
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	err = utils.ValidateUser(&dbUser, user.Password)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	tokenString, err := jwtI.GenerateToken(dbUser)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Could not log in"))
		return
	}

	jwtI.SetTokenCookie(w, tokenString)

	jsonResponse := map[string]string{"token": tokenString}
	utils.SendJSONResponse(w, http.StatusOK, jsonResponse)
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	if r == nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	if rdb == nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	tokenString := r.Header.Get("Authorization")
	actualToken := strings.TrimPrefix(tokenString, "Bearer ")

	if actualToken == "" || actualToken == tokenString {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	token, err := jwtI.ParseToken(actualToken)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	expirationTime := int64(claims["exp"].(float64))

	const maxRetries = 3
	var currentRetry = 0

	for currentRetry < maxRetries {
		err = redisI.BlacklistToken(rdb, tokenString, expirationTime)
		if err == nil {
			break
		}

		currentRetry++
		time.Sleep(2 * time.Second)
	}

	if currentRetry == maxRetries {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Warnf("Failed to blacklist token after %d retries", maxRetries)
	}

	jwtI.ClearTokenCookie(w)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}
