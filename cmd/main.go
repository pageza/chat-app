// Package main is the entry point for the chat application.
// It initializes logging, configuration, database, and starts the server.
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/logging"
	"github.com/pageza/chat-app/internal/server"
	"github.com/pageza/chat-app/pkg/database"
	"github.com/sirupsen/logrus"
)

// main is the entry point function for the chat application.
func main() {

	logrus.Info("Starting logging initialization")
	if err := logging.Initialize(); err != nil {
		logrus.Fatalf("Failed to initialize logging: %v", err)
		return
	}
	logrus.Info("Logging initialized")
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	logrus.Info("Starting config initialization")
	config.Initialize()
	logrus.Info("Config initialized")

	logrus.Info("Starting database initialization")
	database.GetDB()
	logrus.Info("Database initialized")

	logrus.Info("Starting server initialization")
	server.StartServer()
	logrus.Info("Server initialized")

	tokenExpiration := "2h" // This should be loaded from your config
	duration, err := time.ParseDuration(tokenExpiration)
	if err != nil {
		logrus.Fatalf("Invalid token expiration duration: %v", err)
		return
	}
	logrus.Infof("Parsed duration: %v", duration)

	logrus.Info("Application started")
}
