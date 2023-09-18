package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/pageza/chat-app/database"
	"github.com/pageza/chat-app/middleware"
	"github.com/pageza/chat-app/models"

	"github.com/rs/cors"

	"golang.org/x/crypto/bcrypt"
)

// HealthCheckHandler returns the status of the server
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Server is up and running!")
}

// ChatHandler handles chat-related requests
func ChatHandler(w http.ResponseWriter, r *http.Request) {
	// Placeholder for chat functionality
	fmt.Fprintf(w, "Chat endpoint")
}
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Could not hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	result := database.DB.Create(&user)
	if result.Error != nil {
		http.Error(w, "Could not register user", http.StatusInternalServerError)
		return
	}
	// Generate token
	tokenString, err := middleware.GenerateToken(user)
	if err != nil {
		http.Error(w, "Could not log in", http.StatusInternalServerError)
		return
	}
	// Set HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
	})

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User successfully registered and logged in")
}

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dbUser models.User
	database.DB.Where("username = ?", user.Username).First(&dbUser)

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate token and respond (to be implemented)
	// Generate token
	tokenString, err := middleware.GenerateToken(dbUser)
	if err != nil {
		http.Error(w, "Could not log in", http.StatusInternalServerError)
		return
	}
	// Set HttpOnly cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
	})
	// Create a JSON response
	jsonResponse := map[string]string{"token": tokenString}

	// Serialize and send the JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jsonResponse)

}

// LogoutHandler handles user logout
func LogoutHandler(w http.ResponseWriter, r *http.Request, rdb *redis.Client) {
	// Check if request object is nil
	if r == nil {
		http.Error(w, "Request is nil", http.StatusInternalServerError)
		return
	}

	// Check if Redis client is initialized
	if rdb == nil {
		http.Error(w, "Redis client is not initialized", http.StatusInternalServerError)
		return
	}
	// Extracting the token from the request header
	tokenString := r.Header.Get("Authorization")
	fmt.Println("Recieved Token:", tokenString)
	actualToken := strings.TrimPrefix(tokenString, "Bearer ")
	// Parse the token
	token, err := jwt.Parse(actualToken, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})

	if err != nil {
		fmt.Println("Token validation failed: ", err)
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Extract expiration time from the token
	expirationTime := int64(claims["exp"].(float64))

	// Adding the token to the blacklist
	err = rdb.Set(context.TODO(), tokenString, "blacklisted", time.Until(time.Unix(expirationTime, 0))).Err()
	if err != nil {
		fmt.Println("failed to blacklist token", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		MaxAge:   -1, // MaxAge<0 means delete cookie now
	})
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Logged out successfully")
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve username and email from request headers, which were set by ValidateMiddleware
	username := r.Header.Get("username")
	email := r.Header.Get("email")

	// If either username or email is empty, it means the JWT was not validated
	if username == "" || email == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Create a user info JSON response
	userInfo := map[string]string{
		"username": username,
		"email":    email,
	}

	// Convert the map to JSON
	jsonResponse, err := json.Marshal(userInfo)
	if err != nil {
		http.Error(w, "Could not create user info response", http.StatusInternalServerError)
		return
	}

	// Set Content-Type and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// SendMessageHandler handles sending messages
func SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Send message endpoint")
}

// ReceiveMessageHandler handles receiving messages
func ReceiveMessageHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Receive message endpoint")
}
func main() {

	// Initialize CORS middleware
	c := cors.New(cors.Options{
		// AllowedOrigins: []string{"*"}, // Allow all origins
		AllowedOrigins:   []string{"http://localhost:8080"}, // replace with your frontend application's URL
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Add this line
	})

	// Get Redis client
	middleware.InitializeRedis()
	rdb := middleware.GetRedisClient()

	// Initialize the router
	r := mux.NewRouter()

	handler := c.Handler(r)

	// Define routes
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/chat", ChatHandler).Methods("GET")
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		LogoutHandler(w, r, rdb)
	}).Methods("POST")
	r.HandleFunc("/send", SendMessageHandler).Methods("POST")
	r.HandleFunc("/receive", ReceiveMessageHandler).Methods("GET")
	r.HandleFunc("/userinfo", middleware.AuthMiddleware(UserInfoHandler)).Methods("GET")
	r.HandleFunc("/check-auth", middleware.CheckAuth).Methods("GET")

	// Start the HTTP server
	fmt.Println("Server is running on port 8888")
	log.Fatal(http.ListenAndServe(":8888", handler))

}
