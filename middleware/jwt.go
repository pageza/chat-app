package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/pageza/chat-app/models"
)

func GenerateToken(user models.User) (string, error) {
	// Set token expiration time, e.g., 1 hour from now
	expirationTime := time.Now().Add(1 * time.Hour).Unix()

	claims := &jwt.StandardClaims{
		// Set the expiration time
		ExpiresAt: expirationTime,
		// Add other claims
		Issuer:  "veterans-app",
		Subject: user.Username,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("your_secret_key"))

	return tokenString, err
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

		// Parse token
		tokenStr := cookie.Value
		claims := &jwt.StandardClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			// Replace "YourSigningKey" with your actual signing key
			return []byte("your_secret_key"), nil
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

	// // Setting Headers manually for debugging
	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	// w.Header().Set("Access-Control-Allow-Credentials", "true") // This allows cookies

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
		return []byte("your_secret_key"), nil
	})

	if err != nil || !token.Valid {
		log.Printf("Error: %s", err.Error()) // Log error

		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"authenticated": true})
}
