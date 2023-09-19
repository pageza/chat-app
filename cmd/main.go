package main

import (
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/internal/common"
	"github.com/pageza/chat-app/internal/middleware"
	"github.com/pageza/chat-app/internal/routes"
	"github.com/sirupsen/logrus"

	"github.com/rs/cors"
)

func main() {
	// Initialize Logrus and set log file
	logFile, err := os.OpenFile("application.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %s", err)
	}
	defer logFile.Close()
	logrus.SetOutput(logFile)

	err = godotenv.Load()
	if err != nil {
		logrus.Fatal("Error loading .env file")
	}
	common.Initialize()

	// Read environment variables
	corsAllowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	corsAllowedMethods := strings.Split(os.Getenv("CORS_ALLOWED_METHODS"), ",")
	corsAllowedHeaders := strings.Split(os.Getenv("CORS_ALLOWED_HEADERS"), ",")
	serverPort := os.Getenv("SERVER_PORT")

	// Initialize CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{corsAllowedOrigins},
		AllowCredentials: true,
		AllowedMethods:   corsAllowedMethods,
		AllowedHeaders:   corsAllowedHeaders,
	})

	// TODO: Consider moving the Redis initialization to a separate function for better readability
	// Get Redis client
	middleware.InitializeRedis()
	rdb := middleware.GetRedisClient()

	// Initialize the router
	r := mux.NewRouter()
	r.Use(middleware.RecoveryMiddleware)
	routes.InitializeRoutes(r, rdb)
	r.Use(middleware.RateLimitMiddleware)

	handler := c.Handler(r)

	// Start the HTTP server
	logrus.Info("Server is running on port:", serverPort)
	http.ListenAndServe(":"+serverPort, handler)

}
