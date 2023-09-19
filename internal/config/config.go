package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
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
	// Load the .env file
	err := godotenv.Load("/home/zach/projects/chat-app/.env")
	if err != nil {
		logrus.Fatal("Error loading .env file:", err)
	}

	JwtSecret = os.Getenv("JWT_SECRET")
	JwtIssuer = os.Getenv("JWT_ISSUER")
	RedisAddr = os.Getenv("REDIS_ADDR")
	TokenExpiration = os.Getenv("TOKEN_EXPIRATION")
	CorsAllowedOrigins = []string{os.Getenv("CORS_ALLOWED_ORIGINS")}
	CorsAllowedMethods = strings.Split(os.Getenv("CORS_ALLOWED_METHODS"), ",")
	CorsAllowedHeaders = strings.Split(os.Getenv("CORS_ALLOWED_HEADERS"), ",")
	ServerPort = os.Getenv("SERVER_PORT")
	PostgreDSN = os.Getenv("POSTGRE_DSN")

	if JwtSecret == "" || JwtIssuer == "" || RedisAddr == "" || TokenExpiration == "" || PostgreDSN == "" {
		logrus.Fatal("Environment variables are not set")
	}

	_, err = time.ParseDuration(TokenExpiration)
	if err != nil {
		logrus.Fatalf("Invalid token expiration duration: %v", err)
	}
}
