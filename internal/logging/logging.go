// Package logging provides utility functions for setting up and configuring logging in the application.
// It uses the Logrus library for logging and allows for logging to a file.

package logging

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Initialize sets up the logging configurations.
func Initialize() error {
	logDir := "/workspaces/chat-app/logs"

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		err = os.MkdirAll(logDir, 0755)
		if err != nil {
			return err
		}
	}

	logFilePath := logDir + "/application.log"
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetOutput(logFile)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		logrus.Info("Application is shutting down")
		if err := logFile.Close(); err != nil {
			logrus.Errorf("Error closing log file: %s", err)
		}
	}()

	return nil
}
