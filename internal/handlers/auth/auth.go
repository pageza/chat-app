package auth

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/internal/common"
	"github.com/pageza/chat-app/internal/errors" // <-- Updated import
	"github.com/pageza/chat-app/internal/helpers"
	"github.com/pageza/chat-app/internal/models"
	"github.com/pageza/chat-app/pkg/database"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}
	if err := user.Validate(); err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusBadRequest, "Invalid user data"))
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Internal server error"))
		return
	}
	user.Password = string(hashedPassword)

	const maxRetries = 3
	var currentRetry = 0
	var result *gorm.DB // Declare result here

	for currentRetry < maxRetries {
		result = database.DB.Create(&user) // Assign to result
		if result.Error == nil {
			break
		}

		currentRetry++
		time.Sleep(2 * time.Second)
	}

	if currentRetry == maxRetries || result.Error != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Could not register user"))
		return
	}

	tokenString, err := common.GenerateToken(user)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Could not log in"))
		return
	}
	helpers.SetTokenCookie(w, tokenString)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User successfully registered and logged in")
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	// Validate the user
	if err := user.Validate(); err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusBadRequest, "Invalid user data"))
		return
	}

	const maxRetries = 3
	var currentRetry = 0
	var dbUser models.User

	for currentRetry < maxRetries {
		database.DB.Where("username = ?", user.Username).First(&dbUser)
		if dbUser.ID != 0 {
			break
		}

		currentRetry++
		time.Sleep(2 * time.Second)
	}

	if currentRetry == maxRetries || dbUser.ID == 0 {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	err = helpers.ValidateUser(&dbUser, user.Password)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	tokenString, err := common.GenerateToken(dbUser)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Could not log in"))
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
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	if rdb == nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	tokenString := r.Header.Get("Authorization")
	actualToken := strings.TrimPrefix(tokenString, "Bearer ")

	if actualToken == "" || actualToken == tokenString {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	token, err := helpers.ParseToken(actualToken)
	if err != nil {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		errors.RespondWithError(w, errors.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	expirationTime := int64(claims["exp"].(float64))

	const maxRetries = 3
	var currentRetry = 0

	for currentRetry < maxRetries {
		err = helpers.BlacklistToken(rdb, tokenString, expirationTime)
		if err == nil {
			break
		}

		currentRetry++
		time.Sleep(2 * time.Second)
	}

	if currentRetry == maxRetries {
		log.Printf("Failed to blacklist token after %d retries", maxRetries)
	}

	helpers.ClearTokenCookie(w)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}
