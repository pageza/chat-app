// Package redis provides utilities for interacting with Redis.
// This includes initializing the Redis client, blacklisting tokens, and checking rate limits.

package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/internal/config"
	"github.com/sirupsen/logrus"
)

// rdb is the Redis client that will be used throughout the application.
var rdb *redis.Client

// InitializeRedis sets up the Redis client.
func InitializeRedis() {
	// Create a new Redis client
	rdb = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr, // Redis server address
	})

	// Test the Redis connection
	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		// Log fatal error if Redis connection fails
		logrus.Fatalf("Could not connect to Redis: %v", err)
	}
}

// GetRedisClient returns the initialized Redis client.
func GetRedisClient() *redis.Client {
	return rdb
}

// BlacklistToken blacklists a given JWT token in Redis.
func BlacklistToken(rdb *redis.Client, tokenString string, expirationTime int64) error {
	const maxRetries = 3 // Maximum number of retries
	var currentRetry = 0 // Current retry count

	// Retry logic for blacklisting the token
	for currentRetry < maxRetries {
		err := rdb.Set(context.TODO(), tokenString, "blacklisted", time.Until(time.Unix(expirationTime, 0))).Err()
		if err == nil {
			return nil // Operation was successful, return
		}

		// Increment the retry count and wait before the next retry
		currentRetry++
		time.Sleep(2 * time.Second)
	}

	// Log and return an error if maximum retries are reached
	if currentRetry == maxRetries {
		logrus.Printf("Max retries reached, could not blacklist the token")
		return fmt.Errorf("max retries reached, could not blacklist the token")
	}

	return nil // This line is technically unreachable but added for completeness
}

// CheckRateLimit checks the rate limit for a given IP in Redis.
func CheckRateLimit(ip string, rdb *redis.Client) (bool, error) {
	ctx := context.TODO()
	// Increment the request count for this IP
	newCount, err := rdb.Incr(ctx, ip).Result()
	if err != nil {
		// Log the error and return
		logrus.WithFields(logrus.Fields{
			"ip": ip,
		}).Errorf("Redis error: %v", err)
		return false, err
	}

	// Set the key to expire after 1 minute if this is the first request from this IP
	if newCount == 1 {
		if _, err := rdb.Expire(ctx, ip, 1*time.Minute).Result(); err != nil {
			// Log the error and return
			logrus.WithFields(logrus.Fields{
				"ip": ip,
			}).Errorf("Redis error: %v", err)
			return false, err
		}
	}

	// Check if the rate limit has been exceeded (e.g., more than 10 requests per minute)
	if newCount > 10 {
		// Log a warning and return false
		logrus.WithFields(logrus.Fields{
			"ip": ip,
		}).Warn("Rate limit exceeded")
		return false, nil
	}

	return true, nil // Return true if the rate limit has not been exceeded
}
