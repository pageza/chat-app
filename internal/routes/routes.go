// Package routes sets up all the routes for the application.
// This includes routes for health checks, authentication, chat functionality, and user information.
package routes

import (
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/pageza/chat-app/internal/auth"
	"github.com/pageza/chat-app/internal/chat"
	"github.com/pageza/chat-app/internal/middleware"
	"github.com/pageza/chat-app/internal/user"
	"github.com/pageza/chat-app/internal/utils"
	"github.com/pageza/chat-app/pkg/database"
)

func InitializeRoutes(r *mux.Router, rdb *redis.Client, db database.Database) {
	authHandler := &auth.AuthHandler{DB: db}

	// Health check route
	r.HandleFunc("/health", utils.HealthCheckHandler).Methods("GET")

	// Chat-related routes
	r.HandleFunc("/chat", chat.ChatHandler).Methods("GET")
	r.HandleFunc("/send", chat.SendMessageHandler).Methods("POST")
	r.HandleFunc("/receive", chat.ReceiveMessageHandler).Methods("GET")

	// Authentication-related routes
	r.HandleFunc("/register", authHandler.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", authHandler.LoginHandler).Methods("POST")
	// Logout route with inline function to pass Redis client
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		authHandler.LogoutHandler(w, r, rdb)
	}).Methods("POST")

	// User information route with authentication middleware
	r.HandleFunc("/userinfo", middleware.AuthMiddleware(user.UserInfoHandler)).Methods("GET")

	// Route to check if the user is authenticated
	r.HandleFunc("/check-auth", middleware.CheckAuth).Methods("GET")
}
