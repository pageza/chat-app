package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

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

func Initialize() {
	// Load the .env file for sensitive variables
	err := godotenv.Load("/home/zach/projects/chat-app/.env")
	if err != nil {
		logrus.Fatal("Error loading .env file:", err)
	}

	JwtSecret = os.Getenv("JWT_SECRET")
	JwtIssuer = os.Getenv("JWT_ISSUER")
	PostgreDSN = os.Getenv("POSTGRE_DSN")

	if JwtSecret == "" || JwtIssuer == "" || PostgreDSN == "" {
		logrus.Fatal("Sensitive environment variables are not set")
	}
	// Determine the environment
	env := os.Getenv("ENV") // ENV should be either 'development' or 'production'
	if env == "" {
		env = "development" // Default to development
	}

	// Read environment-specific configurations
	dbHost := viper.GetString(fmt.Sprintf("%s.database.host", env))
	dbPort := viper.GetString(fmt.Sprintf("%s.database.port", env))
	dbUsername := viper.GetString(fmt.Sprintf("%s.database.username", env))
	dbPassword := viper.GetString(fmt.Sprintf("%s.database.password", env))
	dbName := viper.GetString(fmt.Sprintf("%s.database.dbname", env))
	dbSSLMode := viper.GetString(fmt.Sprintf("%s.database.sslmode", env))

	// Construct the DSN
	PostgreDSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUsername, dbPassword, dbName, dbSSLMode)

	// Initialize Viper for non-sensitive, environment-specific variables
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err = viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Fatal error config file: %s \n", err)
	}

	RedisAddr = viper.GetString("REDIS_ADDR")
	TokenExpiration = viper.GetString("TOKEN_EXPIRATION")
	CorsAllowedOrigins = []string{viper.GetString("CORS_ALLOWED_ORIGINS")}
	CorsAllowedMethods = strings.Split(viper.GetString("CORS_ALLOWED_METHODS"), ",")
	CorsAllowedHeaders = strings.Split(viper.GetString("CORS_ALLOWED_HEADERS"), ",")
	ServerPort = viper.GetString("SERVER_PORT")

	_, err = time.ParseDuration(TokenExpiration)
	if err != nil {
		logrus.Fatalf("Invalid token expiration duration: %v", err)
	}
}
