package logging

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.Println("Starting TestMain...") // This should appear in your test output

	// Load .env.test
	if err := godotenv.Load("/home/zach/projects/chat-app/.env"); err != nil {
		log.Fatalf("Error loading .env: %v", err)
	}

	// Initialize Viper
	viper.SetConfigFile("/home/zach/projects/chat-app/internal/config/config.yaml") // Adjust the path as needed
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config.yaml: %v", err)
	}

	// Run the tests
	code := m.Run()

	// Teardown if needed

	os.Exit(code)
}
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

// Additional test cases can be added here for different log levels, formats, etc.
