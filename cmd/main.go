package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/internal/common"
	"github.com/pageza/chat-app/internal/middleware"
	"github.com/pageza/chat-app/internal/routes"

	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
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
	fmt.Println("Server is running on port:", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, handler))

}
