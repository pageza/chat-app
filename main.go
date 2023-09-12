package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
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
	})
	// Initialize the router
	r := mux.NewRouter()

	handler := c.Handler(r)

	// Define routes
	r.HandleFunc("/health", HealthCheckHandler).Methods("GET")
	r.HandleFunc("/chat", ChatHandler).Methods("GET")
	r.HandleFunc("/register", RegisterHandler).Methods("POST")
	r.HandleFunc("/login", LoginHandler).Methods("POST")
	r.HandleFunc("/logout", LogoutHandler).Methods("POST")
	r.HandleFunc("/send", SendMessageHandler).Methods("POST")
	r.HandleFunc("/receive", ReceiveMessageHandler).Methods("GET")
	r.HandleFunc("/userinfo", middleware.AuthMiddleware(UserInfoHandler)).Methods("GET")
	r.HandleFunc("/check-auth", middleware.CheckAuth).Methods("GET")

	// Start the HTTP server
	fmt.Println("Server is running on port 8888")
	log.Fatal(http.ListenAndServe(":8888", handler))

}
