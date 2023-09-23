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
	// Using a relative path
	err = godotenv.Load("/home/zach/projects/chat-app/.env") // Adjust the path as needed
	if err != nil {
		logrus.Fatal("Error loading .env file:", err)
	}

	// Set the path for the config file
	// Using a relative path
	viper.SetConfigFile("/home/zach/projects/chat-app/internal/config/config.yaml") // Adjust the path as needed

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file: %s", err)
	}

	// Debugging: Print all settings
	// fmt.Println("All settings:", viper.AllSettings())

	// Read sensitive variables from environment
	JwtSecret = os.Getenv("JWT_SECRET")
	JwtIssuer = os.Getenv("JWT_ISSUER")
	PostgreDSN = os.Getenv("POSTGRE_DSN")
	// fmt.Println("JWT_SECRET:", os.Getenv("JWT_SECRET"))
	// fmt.Println("JWT_ISSUER:", os.Getenv("JWT_ISSUER"))
	// fmt.Println("POSTGRE_DSN:", os.Getenv("POSTGRE_DSN"))

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
	fmt.Println("te: ", TokenExpiration, "TE: ", viper.GetString("TOKEN_EXPIRATION"))
	// Set default value if TokenExpiration is not set
	if TokenExpiration == "" {
		TokenExpiration = "2h" // Default value
	}
	fmt.Println("THIS IS THE EXPIRY", TokenExpiration)
	// New debug statement to print the exact string being passed to ParseDuration
	fmt.Println("Parsing duration from this exact string:", fmt.Sprintf("%q", TokenExpiration))
	fmt.Println("Debug TokenExpiration in test:", TokenExpiration)
	// Validate the token expiration duration
	_, err = time.ParseDuration(TokenExpiration)
	if err != nil {
		fmt.Println(err)
		logrus.Fatalf("Invalid token expiration duration config: %v", err)
	}
}
