package middleware

import (
	"net/http"
	"time"

	"github.com/pageza/chat-app/internal/common"
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		ip := r.RemoteAddr
		lastRequestTime, exists := limiter[ip]

		if exists && time.Since(lastRequestTime) < 1*time.Second {
			// common.RespondWithError(w, "Too many requests", http.StatusTooManyRequests)
			common.RespondWithError(w, common.NewAPIError(http.StatusTooManyRequests, "Too many requests"))
			return
		}

		limiter[ip] = time.Now()
		next.ServeHTTP(w, r)
	})
}
