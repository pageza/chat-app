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
)

func InitializeRoutes(r *mux.Router, rdb *redis.Client) {
	// Define routes
	r.HandleFunc("/health", utils.HealthCheckHandler).Methods("GET")
	r.HandleFunc("/chat", chat.ChatHandler).Methods("GET")
	r.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")
	r.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		auth.LogoutHandler(w, r, rdb)
	}).Methods("POST")
	r.HandleFunc("/send", chat.SendMessageHandler).Methods("POST")
	r.HandleFunc("/receive", chat.ReceiveMessageHandler).Methods("GET")
	r.HandleFunc("/userinfo", middleware.AuthMiddleware(user.UserInfoHandler)).Methods("GET")
	r.HandleFunc("/check-auth", middleware.CheckAuth).Methods("GET")
}
