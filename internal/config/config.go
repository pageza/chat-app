// Package config handles the configuration settings for the chat application.
// It reads sensitive variables from environment files and other configurations from a config file.
package config

import (
	"os"
	"strings"
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
// It reads from both environment variables and configuration files.
func Initialize() {
	// Load the .env file for sensitive variables
	err := godotenv.Load("/home/zach/projects/chat-app/.env")
	if err != nil {
		// Log fatal error if .env file cannot be loaded
		logrus.WithFields(logrus.Fields{
			"file": ".env",
		}).Fatal("Error loading .env file:", err)
	}

	// Read sensitive variables from environment
	JwtSecret = os.Getenv("JWT_SECRET")
	JwtIssuer = os.Getenv("JWT_ISSUER")
	PostgreDSN = os.Getenv("POSTGRE_DSN")

	// Check if sensitive environment variables are set
	if JwtSecret == "" || JwtIssuer == "" || PostgreDSN == "" {
		logrus.WithFields(logrus.Fields{
			"missing_vars": strings.Join([]string{"JWT_SECRET", "JWT_ISSUER", "POSTGRE_DSN"}, ", "),
		}).Fatal("Sensitive environment variables are not set")
	}

	// Determine the environment (development or production)
	env := os.Getenv("ENV")
	if env == "" {
		env = "development" // Default to development if ENV is not set
	}

	// Load non-sensitive configurations from config file
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err = viper.ReadInConfig()
	if err != nil {
		// Log fatal error if config file cannot be read
		logrus.WithFields(logrus.Fields{
			"file": "config",
		}).Fatalf("Fatal error config file: %s \n", err)
	}

	// Read configurations from config file
	RedisAddr = viper.GetString("REDIS_ADDR")
	TokenExpiration = viper.GetString("TOKEN_EXPIRATION")
	CorsAllowedOrigins = []string{viper.GetString("CORS_ALLOWED_ORIGINS")}
	CorsAllowedMethods = strings.Split(viper.GetString("CORS_ALLOWED_METHODS"), ",")
	CorsAllowedHeaders = strings.Split(viper.GetString("CORS_ALLOWED_HEADERS"), ",")
	ServerPort = viper.GetString("SERVER_PORT")

	// Validate the token expiration duration
	_, err = time.ParseDuration(TokenExpiration)
	if err != nil {
		// Log fatal error if token expiration duration is invalid
		logrus.WithFields(logrus.Fields{
			"duration": TokenExpiration,
		}).Fatalf("Invalid token expiration duration: %v", err)
	}
}
