package config

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Configuration variables
var (
	JwtSecret          string
	JwtIssuer          string
	RedisAddr          string
	TokenExpiration    string
	CorsAllowedOrigins []string
	CorsAllowedMethods []string
	CorsAllowedHeaders []string
	ServerPort         string
	PostgreDSN         string
)

// Initialize sets up the application's configuration.
func Initialize() {
	// Debugging: Print the current working directory
	dir, err := os.Getwd()
	if err != nil {
		logrus.Fatal("Getting working directory failed:", err)
	}
	fmt.Println("Current directory:", dir)

	// Load the .env file for sensitive variables
	err = godotenv.Load("/home/zach/projects/chat-app/.env")
	if err != nil {
		logrus.Fatal("Error loading .env file:", err)
	}

	// Set the path for the config file
	viper.SetConfigFile("/home/zach/projects/chat-app/internal/config/config.yaml")

	// Read the config file
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Error reading config file: %s", err)
	}

	// Debugging: Print all settings
	fmt.Println("All settings:", viper.AllSettings())

	// Read sensitive variables from environment
	JwtSecret = os.Getenv("JWT_SECRET")
	JwtIssuer = os.Getenv("JWT_ISSUER")
	PostgreDSN = os.Getenv("POSTGRE_DSN")

	// Check if sensitive environment variables are set
	if JwtSecret == "" || JwtIssuer == "" || PostgreDSN == "" {
		logrus.Fatal("Sensitive environment variables are not set")
	}

	// Read configurations from config file
	RedisAddr = viper.GetString("REDIS_ADDR")
	TokenExpiration = viper.GetString("TOKEN_EXPIRATION")
	CorsAllowedOrigins = viper.GetStringSlice("CORS_ALLOWED_ORIGINS")
	CorsAllowedMethods = viper.GetStringSlice("CORS_ALLOWED_METHODS")
	CorsAllowedHeaders = viper.GetStringSlice("CORS_ALLOWED_HEADERS")
	ServerPort = viper.GetString("SERVER_PORT")
	// Set default value if TokenExpiration is not set
	if TokenExpiration == "" {
		TokenExpiration = "2h" // Default value
	}

	// Validate the token expiration duration
	_, err = time.ParseDuration(TokenExpiration)
	if err != nil {
		logrus.Fatalf("Invalid token expiration duration: %v", err)
	}
}
