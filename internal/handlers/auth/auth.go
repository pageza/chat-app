package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/internal/common"
	"github.com/pageza/chat-app/internal/helpers"
	"github.com/pageza/chat-app/internal/models"
	"github.com/pageza/chat-app/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}
	// Validate the user
	if err := user.Validate(); err != nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusBadRequest, "Invalid user data"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusInternalServerError, "Internal server error"))
		return
	}
	user.Password = string(hashedPassword)

	result := database.DB.Create(&user)
	if result.Error != nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusInternalServerError, "Could not register user"))
		return
	}

	tokenString, err := common.GenerateToken(user)
	if err != nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusInternalServerError, "Could not log in"))
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
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusBadRequest, "Invalid request payload"))
		return
	}

	// Validate the user
	if err := user.Validate(); err != nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusBadRequest, "Invalid user data"))
		return
	}

	var dbUser models.User
	database.DB.Where("username = ?", user.Username).First(&dbUser)

	err = helpers.ValidateUser(&dbUser, user.Password)
	if err != nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	tokenString, err := common.GenerateToken(dbUser)
	if err != nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusInternalServerError, "Could not log in"))
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
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	if rdb == nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	tokenString := r.Header.Get("Authorization")
	actualToken := strings.TrimPrefix(tokenString, "Bearer ")

	if actualToken == "" || actualToken == tokenString {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	token, err := helpers.ParseToken(actualToken)
	if err != nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusUnauthorized, "Invalid credentials"))
		return
	}

	expirationTime := int64(claims["exp"].(float64))
	err = helpers.BlacklistToken(rdb, tokenString, expirationTime)
	if err != nil {
		helpers.RespondWithError(w, helpers.NewAPIError(http.StatusInternalServerError, "Internal server error"))
		return
	}

	helpers.ClearTokenCookie(w)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}
