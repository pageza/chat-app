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

func StartServer() {
	// Initialize CORS
	c := config.InitializeCORS()

	// Initialize Redis
	rdb := redis.GetRedisClient()

	// Initialize the router
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.RateLimitMiddleware)

	// Initialize Routes
	routes.InitializeRoutes(r, rdb)

	// CORS Handler
	handler := c.Handler(r)

	// Start the HTTP server
	serverPort := os.Getenv("SERVER_PORT")
	logrus.Info("Server is running on port:", serverPort)
	http.ListenAndServe(":"+serverPort, handler)
}
