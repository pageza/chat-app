package middleware

import (
	"net/http"

	"github.com/pageza/chat-app/internal/errors"
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
