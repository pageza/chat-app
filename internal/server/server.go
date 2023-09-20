// Package server initializes and starts the HTTP server.
// It sets up CORS, middleware, and routes.

package server

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pageza/chat-app/internal/config" // Import the config package
	"github.com/pageza/chat-app/internal/middleware"
	"github.com/pageza/chat-app/internal/redis"
	"github.com/pageza/chat-app/internal/routes"
	"github.com/sirupsen/logrus"
)

// StartServer initializes and starts the HTTP server.
func StartServer() {
	// Initialize CORS settings from the config package
	c := config.InitializeCORS()

	// Get the Redis client instance
	rdb := redis.GetRedisClient()

	// Initialize the Gorilla Mux router
	r := mux.NewRouter()

	// Add middleware for recovery and rate-limiting
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.RateLimitMiddleware)

	// Initialize all routes for the application
	routes.InitializeRoutes(r, rdb)

	// Apply CORS settings to the router
	handler := c.Handler(r)

	// Get the server port from environment variables
	serverPort := os.Getenv("SERVER_PORT")

	// Log the server start event
	logrus.Info("Server is running on port:", serverPort)

	// Start the HTTP server
	http.ListenAndServe(":"+serverPort, handler)
}
