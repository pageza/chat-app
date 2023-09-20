// Package logging provides utility functions for setting up and configuring logging in the application.
// It uses the Logrus library for logging and allows for logging to a file.

package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Initialize sets up the logging configurations.
// It creates a log directory if it doesn't exist and sets up a log file for the application.
// Logrus library is configured to write logs to this file.
func Initialize() {
	// Define the log directory. Replace with your actual absolute path.
	// Use an absolute path for the log directory or make sure the relative path is correct.
	logDir := "/home/zach/projects/chat-app/logs"

	// Check if the log directory exists, if not create it.
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		// Use MkdirAll to create nested directories if needed.
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			// Log and terminate if directory creation fails.
			logrus.Fatalf("Failed to create log directory: %s", err)
		}
	}

	// Define the path for the log file.
	logFilePath := logDir + "/application.log"

	// Open or create the log file.
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		// Log and terminate if log file operation fails.
		logrus.Fatalf("Failed to open log file: %s", err)
	}
	defer logFile.Close()

	// Configure Logrus settings.
	// Set the log level to Debug.
	logrus.SetLevel(logrus.DebugLevel)
	// Set the output destination for the logs.
	logrus.SetOutput(logFile)
}
