// Package auth provides authentication handlers for the chat application.
// It includes functionalities for user registration, login, and logout.
package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/internal/errors"
	jwtI "github.com/pageza/chat-app/internal/jwt"
	"github.com/pageza/chat-app/internal/models"
	redisI "github.com/pageza/chat-app/internal/redis"
	"github.com/pageza/chat-app/internal/utils"
	"github.com/pageza/chat-app/pkg/database"
	"github.com/sirupsen/logrus"
)

type AuthHandler struct {
	DB         database.Database
	JwtManager jwtI.JwtManager // Add this line
}

// RedisClient is an interface representing the methods of the Redis client
// that are used in this package.
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	// Add other methods as needed
}

func (a *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entering RegisterHandler")
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

	err = a.DB.CreateUser(&user)
	fmt.Println("Debug: CreateUser error:", err)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Warn("Could not register user after max retries")
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Could not register user"))
		return
	}

	accessToken, refreshToken, err := a.JwtManager.GenerateToken(user)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Could not log in"))
		return
	}
	a.JwtManager.SetTokenCookie(w, accessToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User successfully registered and logged in")
}

func (a *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Debug: Starting LoginHandler") // Debug print

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		fmt.Println("Debug: Invalid request payload") // Debug print
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Warn("Invalid request payload")
		errors.RespondWithError(w, errors.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	dbUser, err := a.DB.GetUserByUsername(user.Username)
	if err != nil || dbUser == nil {
		fmt.Println("Debug: User not found") // Debug print
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Warn("User not found")
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "User not found"))
		return
	}

	err = utils.ValidateUser(dbUser, user.Password)
	if err != nil {
		fmt.Println("Debug: Invalid password") // Debug print
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Warn("Invalid password")
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid password"))
		return
	}

	fmt.Println("Debug: About to call GenerateToken") // Debug print
	fmt.Printf("Debug: dbUser type: %T, content: %+v\n", dbUser, dbUser)

	accessToken, refreshToken, err := a.JwtManager.GenerateToken(*dbUser)

	fmt.Println("Debug: GenerateToken called") // Debug print

	if err != nil {
		fmt.Println("Debug: Entering error block")      // Debug print
		fmt.Println("Debug: GenerateToken error:", err) // Debug print
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Could not log in"))
		return
	}

	a.JwtManager.SetTokenCookie(w, accessToken)

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
	})

	jsonResponse := map[string]string{"token": accessToken}
	utils.SendJSONResponse(w, http.StatusOK, jsonResponse)
}

// The rest of the file remains the same.

// LogoutHandler handles user logout.
// It invalidates the user's JWT token and removes it from the cookie.
func (a *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request, redisClient RedisClient) {
	fmt.Println("Entering LogoutHandler")
	// Validate the incoming request and Redis client
	if r == nil || redisClient == nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	// Extract and validate the JWT token from the Authorization header
	tokenString := r.Header.Get("Authorization")
	actualToken := strings.TrimPrefix(tokenString, "Bearer ")
	if actualToken == "" || actualToken == tokenString {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	// Parse and validate the JWT token
	token, err := a.JwtManager.ParseToken(actualToken)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	// Blacklist the JWT token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}
	expirationTime := int64(claims["exp"].(float64))
	const maxRetries = 3
	var currentRetry = 0
	for currentRetry < maxRetries {
		err = redisI.BlacklistToken(context.TODO(), redisClient, actualToken, expirationTime)
		if err == nil {
			break
		}
		currentRetry++
		time.Sleep(2 * time.Second)
	}

	// Handle logout failure
	if currentRetry == maxRetries {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
			"ip":     r.RemoteAddr,
		}).Warnf("Failed to blacklist token after %d retries", maxRetries)
	}

	// Clear the JWT cookie and send a success response
	a.JwtManager.ClearTokenCookie(w)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}
