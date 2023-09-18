package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/models"
)

var (
	jwtSecret       string
	jwtIssuer       string
	redisAddr       string
	tokenExpiration string
)

// In middleware package
func Initialize() {
	jwtSecret = os.Getenv("JWT_SECRET")
	jwtIssuer = os.Getenv("JWT_ISSUER")
	redisAddr = os.Getenv("REDIS_ADDR")
	tokenExpiration = os.Getenv("TOKEN_EXPIRATION")

	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set in the environment")
	}
	if jwtIssuer == "" {
		log.Fatal("JWT_ISSUER is not set in the environment")
	}
	if redisAddr == "" {
		log.Fatal("REDIS_ADDR is not set in the environment")
	}
	if tokenExpiration == "" {
		log.Fatal("TOKEN_EXPIRATION is not set in the environment")
	}
	_, err := time.ParseDuration(tokenExpiration)
	if err != nil {
		log.Fatalf("Invalid token expiration duration: %v", err)
	}
}

func GenerateToken(user models.User) (string, error) {
	// Convert TOKEN_EXPIRATION to time.Duration
	expirationDuration, err := time.ParseDuration(tokenExpiration)
	if err != nil {
		log.Fatalf("Invalid token expiration duration: %v", err)
	}
	// Set token expiration time, e.g., 1 hour from now
	expirationTime := time.Now().Add(expirationDuration).Unix()

	claims := &jwt.StandardClaims{
		// Set the expiration time
		ExpiresAt: expirationTime,
		// Add other claims
		Issuer:  jwtIssuer,
		Subject: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))

	return tokenString, err
}

// Initializing Redis for token blacklisting
var rdb *redis.Client

func InitializeRedis() {

	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	// Ping the Redis server to check if it's up
	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	} else {
		log.Println("Successfully connected to Redis")
	}
}
func GetRedisClient() *redis.Client {
	return rdb
}

// AuthMiddleware checks for a valid JWT token in the HttpOnly cookie
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Missing auth token", http.StatusUnauthorized)
			return
		}
		// Check if token is blacklisted
		isBlacklisted, err := rdb.Get(context.TODO(), cookie.Value).Result()
		if err == nil && isBlacklisted == "blacklisted" {
			http.Error(w, "Token is blacklisted", http.StatusUnauthorized)
			return
		}
		// Parse token
		tokenStr := cookie.Value
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			// Replace "YourSigningKey" with your actual signing key
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid auth token", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

// CheckAuth checks if a user is logged in and responds with a JSON object.
func CheckAuth(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for CheckAuth from %s", r.RemoteAddr) // Log incoming request

	// Handle preflight request. Needed for CORS support to work.
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Get token from cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	// Parse token
	tokenStr := cookie.Value
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// Replace "YourSigningKey" with your actual signing key
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Error: %s", err.Error()) // Log error

		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"authenticated": true})
}
