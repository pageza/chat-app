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
	c := config.InitializeCORS() // Use the InitializeCORS function from the config package

	// Initialize Redis
	rdb := redis.GetRedisClient() // Initialize rdb here

	// Initialize the router
	r := mux.NewRouter()
	r.Use(middleware.RecoveryMiddleware)
	routes.InitializeRoutes(r, rdb) // Now rdb is defined
	r.Use(middleware.RateLimitMiddleware)

	handler := c.Handler(r)

	// Start the HTTP server
	serverPort := os.Getenv("SERVER_PORT")
	logrus.Info("Server is running on port:", serverPort)
	http.ListenAndServe(":"+serverPort, handler)
}
