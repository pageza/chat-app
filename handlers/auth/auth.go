package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/database"
	"github.com/pageza/chat-app/middleware"
	"github.com/pageza/chat-app/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		handleError(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	result := database.DB.Create(&user)
	if result.Error != nil {
		handleError(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	tokenString, err := middleware.GenerateToken(user)
	if err != nil {
		handleError(w, "Could not log in", http.StatusInternalServerError)
		return
	}

	setTokenCookie(w, tokenString)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User successfully registered and logged in")
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		handleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dbUser models.User
	database.DB.Where("username = ?", user.Username).First(&dbUser)

	err = validateUser(&dbUser, user.Password)
	if err != nil {
		handleError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	tokenString, err := middleware.GenerateToken(dbUser)
	if err != nil {
		handleError(w, "Could not log in", http.StatusInternalServerError)
		return
	}

	setTokenCookie(w, tokenString)

	jsonResponse := map[string]string{"token": tokenString}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	if r == nil {
		handleError(w, "Request is nil", http.StatusInternalServerError)
		return
	}

	if rdb == nil {
		handleError(w, "Redis client is not initialized", http.StatusInternalServerError)
		return
	}

	tokenString := r.Header.Get("Authorization")
	actualToken := strings.TrimPrefix(tokenString, "Bearer ")

	token, err := parseToken(actualToken)
	if err != nil {
		handleError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		handleError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	expirationTime := int64(claims["exp"].(float64))
	err = blacklistToken(rdb, tokenString, expirationTime)
	if err != nil {
		handleError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1,
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}

func validateUser(user *models.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err
}
func setTokenCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	})
}
func handleError(w http.ResponseWriter, message string, statusCode int) {
	http.Error(w, message, statusCode)
}
func parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})
}
func blacklistToken(rdb *redis.Client, tokenString string, expirationTime int64) error {
	return rdb.Set(context.TODO(), tokenString, "blacklisted", time.Until(time.Unix(expirationTime, 0))).Err()
}
