// Package middleware provides utility functions for handling middleware logic in the application.
// This file specifically includes a middleware for recovering from panics.

package middleware

import (
	"net/http"

	"github.com/pageza/chat-app/internal/errors"
	"github.com/sirupsen/logrus"
)

// RecoveryMiddleware is a middleware function that recovers from panics in the application.
// If a panic occurs, it logs the error and returns a 500 Internal Server Error response.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use defer to ensure the function call is executed later in case of a panic
		defer func() {
			// Recover from panic and log the error
			if err := recover(); err != nil {
				logrus.WithFields(logrus.Fields{
					"method": r.Method,
					"url":    r.URL.String(),
					"ip":     r.RemoteAddr,
				}).Errorf("Recovered from panic: %v", err)
				// Respond with a 500 Internal Server Error
				errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Internal Server Error"))
			}
		}()
		// Call the next middleware or handler
		next.ServeHTTP(w, r)
	})
}
