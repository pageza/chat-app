// Package main is the entry point for the chat application.
// It initializes logging, configuration, database, and starts the server.
package main

import (
	"github.com/pageza/chat-app/internal/config"
	"github.com/pageza/chat-app/internal/logging"
	"github.com/pageza/chat-app/internal/server"
	"github.com/pageza/chat-app/pkg/database"
	"github.com/sirupsen/logrus"
)

// main is the entry point function for the chat application.
func main() {

	// Initialize logging: sets up the logging configurations.
	logging.Initialize()

	// Initialize configuration: loads environment variables and sets up configurations.
	config.Initialize()

	// Initialize the database: connects to the database and returns a handle to it.
	database.GetDB()

	// Start the server: initializes routes and starts listening for incoming HTTP requests.
	server.StartServer()

	// Log that the application has started: useful for debugging and monitoring.
	logrus.Info("Application started")

}
