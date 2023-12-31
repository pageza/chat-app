// Package main is the entry point for the chat application.
// It initializes logging, configuration, database, and starts the server.
package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

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
	logrus.Infof("Config - TokenExpiration: %s", config.TokenExpiration)
	logrus.Info("Config initialized")

	logrus.Info("Starting database initialization")
	db, err := database.NewGormDatabase()
	if err != nil {
		logrus.Fatalf("Database initialization failed: %v", err)
		return
	}

	logrus.Info("Database initialized")

	logrus.Info("Starting server initialization")
	server.StartServer(db)
	logrus.Info("Server initialized")

	logrus.Info("Application started")
}
