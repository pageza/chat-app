package config

import (
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
		logrus.WithFields(logrus.Fields{
			"file": ".env",
		}).Fatal("Error loading .env file:", err)
	}

	JwtSecret = os.Getenv("JWT_SECRET")
	JwtIssuer = os.Getenv("JWT_ISSUER")
	PostgreDSN = os.Getenv("POSTGRE_DSN")

	if JwtSecret == "" || JwtIssuer == "" || PostgreDSN == "" {
		logrus.WithFields(logrus.Fields{
			"missing_vars": strings.Join([]string{"JWT_SECRET", "JWT_ISSUER", "POSTGRE_DSN"}, ", "),
		}).Fatal("Sensitive environment variables are not set")
	}

	// Determine the environment
	env := os.Getenv("ENV") // ENV should be either 'development' or 'production'
	if env == "" {
		env = "development" // Default to development
	}

	// Read environment-specific configurations
	viper.SetConfigName("config")
	viper.AddConfigPath("./config/")
	err = viper.ReadInConfig()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"file": "config",
		}).Fatalf("Fatal error config file: %s \n", err)
	}

	RedisAddr = viper.GetString("REDIS_ADDR")
	TokenExpiration = viper.GetString("TOKEN_EXPIRATION")
	CorsAllowedOrigins = []string{viper.GetString("CORS_ALLOWED_ORIGINS")}
	CorsAllowedMethods = strings.Split(viper.GetString("CORS_ALLOWED_METHODS"), ",")
	CorsAllowedHeaders = strings.Split(viper.GetString("CORS_ALLOWED_HEADERS"), ",")
	ServerPort = viper.GetString("SERVER_PORT")

	_, err = time.ParseDuration(TokenExpiration)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"duration": TokenExpiration,
		}).Fatalf("Invalid token expiration duration: %v", err)
	}
}
