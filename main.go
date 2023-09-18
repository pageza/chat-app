package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/middleware"
	"github.com/pageza/chat-app/routes"

	"github.com/rs/cors"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	middleware.Initialize()

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

	// Get Redis client
	middleware.InitializeRedis()
	rdb := middleware.GetRedisClient()

	// Initialize the router
	r := mux.NewRouter()
	routes.InitializeRoutes(r, rdb)

	handler := c.Handler(r)

	// Start the HTTP server
	fmt.Println("Server is running on port:", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, handler))

}
