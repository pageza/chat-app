package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/pageza/chat-app/models"
)

func GenerateToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"email":    user.Email,
	})

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
			return []byte("YourSigningKey"), nil
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
		return []byte("YourSigningKey"), nil
	})

	if err != nil || !token.Valid {
		json.NewEncoder(w).Encode(map[string]bool{"authenticated": false})
		return
	}

	json.NewEncoder(w).Encode(map[string]bool{"authenticated": true})
}

// Old way of making sure its the right user
// func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authorizationHeader := r.Header.Get("Authorization")
// 		if authorizationHeader == "" {
// 			http.Error(w, "Authorization header required", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenString := strings.Split(authorizationHeader, " ")[1]
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			return []byte("your_secret_key"), nil
// 		})

// 		if err != nil {
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 			return
// 		}

// 		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
// 			r.Header.Set("username", claims["username"].(string))
// 			r.Header.Set("email", claims["email"].(string))
// 			next.ServeHTTP(w, r)
// 		} else {
// 			http.Error(w, "Invalid token", http.StatusUnauthorized)
// 		}
// 	})
// }
