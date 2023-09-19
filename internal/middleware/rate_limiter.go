package middleware

import (
	"net/http"

	"github.com/pageza/chat-app/internal/redis"
	"github.com/sirupsen/logrus"
)

// RateLimitMiddleware is a middleware for rate limiting based on IP address
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		userAgent := r.UserAgent()
		rdb := redis.GetRedisClient()

		allowed, err := redis.CheckRateLimit(ip, rdb)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"ip":        ip,
				"userAgent": userAgent,
			}).Errorf("Rate limit check failed: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			logrus.WithFields(logrus.Fields{
				"ip":        ip,
				"userAgent": userAgent,
			}).Warn("Rate limit exceeded")
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
