package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pageza/chat-app/internal/common"
	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/redis"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			// common.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
			common.RespondWithError(w, common.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
			return
		}
		rdb := redis.GetRedisClient()
		isBlacklisted, err := rdb.Get(context.TODO(), cookie.Value).Result()
		if err == nil && isBlacklisted == "blacklisted" {
			// common.RespondWithError(w, "Unauthorized", http.StatusUnauthorized)
			common.RespondWithError(w, common.NewAPIError(http.StatusUnauthorized, "Unauthorized"))
			return
		}

		tokenStr := cookie.Value
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JwtSecret), nil
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
		return []byte(config.JwtSecret), nil
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
