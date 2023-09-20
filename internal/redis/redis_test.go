package redis

import (
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

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
