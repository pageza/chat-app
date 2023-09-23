package redis

import (
	"log"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
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

func TestCheckRateLimit(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Replace with your Redis server address
	})

	ip := "192.168.1.1"

	// Simulate 9 requests
	for i := 0; i < 9; i++ {
		allowed, err := CheckRateLimit(ip, rdb)
		assert.Nil(t, err)
		assert.True(t, allowed)
	}

	// Simulate the 10th request, should still be allowed
	allowed, err := CheckRateLimit(ip, rdb)
	assert.Nil(t, err)
	assert.True(t, allowed)

	// Simulate the 11th request, should be rate-limited
	allowed, err = CheckRateLimit(ip, rdb)
	assert.Nil(t, err)
	assert.False(t, allowed)
}
