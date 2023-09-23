package config_test

import (
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/pageza/chat-app/internal/config"
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

// TestInitialize tests the Initialize function in the config package.
func TestInitialize(t *testing.T) {
	// Backup original environment variables and defer their restoration
	originalJwtSecret := os.Getenv("JWT_SECRET")
	originalJwtIssuer := os.Getenv("JWT_ISSUER")
	originalPostgreDSN := os.Getenv("POSTGRE_DSN")
	defer func() {
		os.Setenv("JWT_SECRET", originalJwtSecret)
		os.Setenv("JWT_ISSUER", originalJwtIssuer)
		os.Setenv("POSTGRE_DSN", originalPostgreDSN)
	}()

	// Set environment variables for testing
	os.Setenv("JWT_SECRET", "test_secret")
	os.Setenv("JWT_ISSUER", "test_issuer")
	os.Setenv("POSTGRE_DSN", "test_dsn")

	// Initialize config
	config.Initialize()

	// Assert that the environment variables were correctly loaded into config variables
	assert.Equal(t, "test_secret", config.JwtSecret, "JwtSecret should be set to the value of the JWT_SECRET environment variable")
	assert.Equal(t, "test_issuer", config.JwtIssuer, "JwtIssuer should be set to the value of the JWT_ISSUER environment variable")
	assert.Equal(t, "test_dsn", config.PostgreDSN, "PostgreDSN should be set to the value of the POSTGRE_DSN environment variable")
}
