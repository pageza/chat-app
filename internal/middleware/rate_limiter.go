// Package middleware provides utility functions for handling middleware logic in the application.
// This file specifically includes a middleware for rate limiting based on the IP address.

package middleware

import (
	"net/http"

	"github.com/pageza/chat-app/internal/redis"
	"github.com/sirupsen/logrus"
)

// RateLimitMiddleware is a middleware function that enforces rate limiting based on the client's IP address.
// It uses Redis to keep track of the number of requests from each IP and blocks any IP that exceeds the limit.
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the client's IP address and User-Agent
		ip := r.RemoteAddr
		userAgent := r.UserAgent()

		// Get the Redis client
		rdb := redis.GetRedisClient()

		// Check if the IP is allowed to make a request
		allowed, err := redis.CheckRateLimit(ip, rdb)
		if err != nil {
			// Log the error and respond with a 500 Internal Server Error
			logrus.WithFields(logrus.Fields{
				"ip":        ip,
				"userAgent": userAgent,
			}).Errorf("Rate limit check failed: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// If the IP has exceeded the rate limit, respond with a 429 Too Many Requests
		if !allowed {
			logrus.WithFields(logrus.Fields{
				"ip":        ip,
				"userAgent": userAgent,
			}).Warn("Rate limit exceeded")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Call the next middleware or handler
		next.ServeHTTP(w, r)
	})
}
