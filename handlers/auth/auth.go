package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/database"
	"github.com/pageza/chat-app/helpers"
	"github.com/pageza/chat-app/middleware"
	"github.com/pageza/chat-app/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Validate the user
	if err := user.Validate(); err != nil {
		helpers.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.HandleError(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	result := database.DB.Create(&user)
	if result.Error != nil {
		helpers.HandleError(w, "Could not register user", http.StatusInternalServerError)
		return
	}

	tokenString, err := middleware.GenerateToken(user)
	if err != nil {
		helpers.HandleError(w, "Could not log in", http.StatusInternalServerError)
		return
	}

	helpers.SetTokenCookie(w, tokenString)

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User successfully registered and logged in")
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate the user
	if err := user.Validate(); err != nil {
		helpers.HandleError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dbUser models.User
	database.DB.Where("username = ?", user.Username).First(&dbUser)

	err = helpers.ValidateUser(&dbUser, user.Password)
	if err != nil {
		helpers.HandleError(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	tokenString, err := middleware.GenerateToken(dbUser)
	if err != nil {
		helpers.HandleError(w, "Could not log in", http.StatusInternalServerError)
		return
	}

	helpers.SetTokenCookie(w, tokenString)

	// Creating the JSON response
	jsonResponse := map[string]string{"token": tokenString}

	// Serializing and sending the JSON response using helper function
	helpers.SendJSONResponse(w, http.StatusOK, jsonResponse)
}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	if r == nil {
		helpers.HandleError(w, "Request is nil", http.StatusInternalServerError)
		return
	}

	if rdb == nil {
		helpers.HandleError(w, "Redis client is not initialized", http.StatusInternalServerError)
		return
	}

	tokenString := r.Header.Get("Authorization")
	actualToken := strings.TrimPrefix(tokenString, "Bearer ")

	if actualToken == "" || actualToken == tokenString {
		helpers.HandleError(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	token, err := helpers.ParseToken(actualToken)
	if err != nil {
		helpers.HandleError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		helpers.HandleError(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	expirationTime := int64(claims["exp"].(float64))
	err = helpers.BlacklistToken(rdb, tokenString, expirationTime)
	if err != nil {
		helpers.HandleError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	helpers.ClearTokenCookie(w)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}
