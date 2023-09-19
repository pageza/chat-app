package main

import (
	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/logging"
	"github.com/pageza/chat-app/internal/server"
	"github.com/pageza/chat-app/pkg/database"
	"github.com/sirupsen/logrus"
)

func main() {

	// Initialize logging
	logging.Initialize()

	// Initialize configuration
	config.Initialize()

	// Initialize the database
	database.GetDB()

	// Start the server
	server.StartServer()

	// Log that the application has started
	logrus.Info("Application started")

}
