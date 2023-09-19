package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/internal/common"
)

var rdb *redis.Client

func InitializeRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: common.RedisAddr,
	})

	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
}

func GetRedisClient() *redis.Client {
	return rdb
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			// common.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
			common.RespondWithError(w, common.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
			return
		}

		isBlacklisted, err := rdb.Get(context.TODO(), cookie.Value).Result()
		if err == nil && isBlacklisted == "blacklisted" {
			// common.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
			common.RespondWithError(w, common.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
			return
		}

		tokenStr := cookie.Value
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(common.JWT_SECRET), nil
		})

		if err != nil || !token.Valid {
			// common.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
			common.RespondWithError(w, common.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
			return
		}

		next(w, r)
	})
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		// common.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		common.RespondWithError(w, common.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	tokenStr := cookie.Value
	claims := &jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(common.JWT_SECRET), nil
	})

	if err != nil || !token.Valid {
		// common.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
		common.RespondWithError(w, common.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"authenticated": true})
}

var (
	limiter = make(map[string]time.Time)
	mu      sync.Mutex
)

func RateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		ip := r.RemoteAddr
		lastRequestTime, exists := limiter[ip]

		if exists && time.Since(lastRequestTime) < 1*time.Second {
			// common.RespondWithError(w, "Too many requests", http.StatusTooManyRequests)
			common.RespondWithError(w, common.NewAPIError(http.StatusTooManyRequests, "Too many requests"))
			return
		}

		limiter[ip] = time.Now()
		next.ServeHTTP(w, r)
	})
}
