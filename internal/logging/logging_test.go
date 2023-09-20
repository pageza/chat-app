package logging

import (
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	// Temporarily set logDir for testing
	logDir := "./temp_logs"
	logFilePath := logDir + "/application.log"

	// Cleanup any existing test log directory
	os.RemoveAll(logDir)

	// Initialize logging
	Initialize()

	// Check if log directory exists
	_, err := os.Stat(logDir)
	assert.False(t, os.IsNotExist(err), "Log directory should exist")

	// Check if log file exists
	_, err = os.Stat(logFilePath)
	assert.False(t, os.IsNotExist(err), "Log file should exist")

	// Check Logrus configuration
	assert.Equal(t, logrus.DebugLevel, logrus.GetLevel(), "Logrus level should be set to Debug")

	// Cleanup test log directory
	os.RemoveAll(logDir)
}
