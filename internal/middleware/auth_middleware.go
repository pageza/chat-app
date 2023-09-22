// Package middleware provides utility functions for handling middleware logic in the application.
// It includes authentication and rate-limiting middleware.

package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pageza/chat-app/internal/common"
	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/redis"
	"github.com/sirupsen/logrus"
)

// AuthMiddleware is a middleware function for handling authentication.
// It checks for a valid JWT token in the request cookie and proceeds to the next handler if valid.
// AuthMiddleware is a middleware function for handling authentication.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !validateToken(r) {
			unauthorizedAccess(w, r)
			return
		}
		next(w, r)
	})
}

// validateToken validates the JWT token from the request.
func validateToken(r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil {
		return false
	}

	// Check if the token is blacklisted
	rdb := redis.GetRedisClient()
	isBlacklisted, err := rdb.Get(context.TODO(), cookie.Value).Result()
	if err == nil && isBlacklisted == "blacklisted" {
		return false
	}

	tokenStr := cookie.Value
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecret), nil
	})

	return err == nil && token.Valid
}

// unauthorizedAccess logs and responds to unauthorized access attempts.
func unauthorizedAccess(w http.ResponseWriter, r *http.Request) {
	logrus.WithFields(logrus.Fields{
		"method": r.Method,
		"url":    r.URL.String(),
		"ip":     r.RemoteAddr,
	}).Warn("Unauthorized access attempt")
	common.RespondWithError(w, common.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
}

// CheckAuth is a utility function to check if the request is authenticated.
// It checks for a valid JWT token in the request cookie and responds with the authentication status.
func CheckAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Retrieve and validate the token from the cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		common.RespondWithError(w, common.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	tokenStr := cookie.Value
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.JwtSecret), nil
	})

	if err != nil || !token.Valid {
		common.RespondWithError(w, common.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	// Respond with the authentication status
	json.NewEncoder(w).Encode(map[string]bool{"authenticated": true})
}

// limiter is a map to store the timestamp of the last request for each IP address.
var limiter = make(map[string]time.Time)

// mu is a mutex for synchronizing access to the limiter map.
var mu sync.Mutex

// ... (Rate limiting logic can go here)
