package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/internal/models"
)

var (
	jwtSecret       string
	jwtIssuer       string
	redisAddr       string
	tokenExpiration string
)

// TODO: Consider adding middleware for logging and metrics.

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
		cookie, err := r.Cookie("token")
		if err != nil {
			log.Println("Missing auth token")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		isBlacklisted, err := rdb.Get(context.TODO(), cookie.Value).Result()
		if err == nil && isBlacklisted == "blacklisted" {
			log.Println("Token is blacklisted")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			log.Printf("Invalid auth token: %s", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	})
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request for CheckAuth from %s", r.RemoteAddr)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		log.Printf("Error: %s", err.Error())
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	tokenStr := cookie.Value
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Error: %s", err.Error())
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"authenticated": true})
}

// Implementing Rate limiting
var (
	limiter = make(map[string]time.Time)
	mu      sync.Mutex
)

// RateLimitMiddleware rate limits requests
func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		ip := r.RemoteAddr
		lastRequestTime, exists := limiter[ip]

		if exists && time.Since(lastRequestTime) < 1*time.Second {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		limiter[ip] = time.Now()
		next.ServeHTTP(w, r)
	})
}
