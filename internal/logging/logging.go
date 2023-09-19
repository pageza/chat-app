package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Initialize sets up the logging configurations
func Initialize() {
	// Use an absolute path for the log directory or make sure the relative path is correct
	logDir := "/home/zach/projects/chat-app/logs" // Replace with your actual absolute path
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755) // Use MkdirAll to create nested directories if needed
		if err != nil {
			logrus.Fatalf("Failed to create log directory: %s", err)
		}
	}

	logFilePath := logDir + "/application.log"
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %s", err)
	}
	defer logFile.Close()

	// Set Logrus configurations
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(logFile)
}
