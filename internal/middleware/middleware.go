package middleware

import (
	"net/http"

	"github.com/pageza/chat-app/internal/errors"
	"github.com/pageza/chat-app/internal/redis"
	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware recovers from panics and writes a 500 if anything goes wrong.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Printf("Recovered from panic: %v", err)
				errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Internal Server Error"))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// RateLimitMiddleware is a middleware for rate limiting based on IP address
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		rdb := redis.GetRedisClient() // Assuming you have a GetRedisClient function in your redis package

		allowed, err := redis.CheckRateLimit(ip, rdb)
		if err != nil {
			logrus.Errorf("Rate limit check failed: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
