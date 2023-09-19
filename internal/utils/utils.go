package utils

import (
	"fmt"
	"net/http"

	"github.com/pageza/chat-app/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// HealthCheckHandler returns the status of the server
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Server is up and running!")
}

func ValidateUser(user *models.User, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err
}
