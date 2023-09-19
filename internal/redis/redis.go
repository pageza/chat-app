package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/pageza/chat-app/internal/config"
	"github.com/sirupsen/logrus"
)

var rdb *redis.Client

func InitializeRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr: config.RedisAddr,
	})

	_, err := rdb.Ping(context.TODO()).Result()
	if err != nil {
		logrus.Fatalf("Could not connect to Redis: %v", err)
	}
}

func GetRedisClient() *redis.Client {
	return rdb
}

func BlacklistToken(rdb *redis.Client, tokenString string, expirationTime int64) error {
	const maxRetries = 3
	var currentRetry = 0

	for currentRetry < maxRetries {
		err := rdb.Set(context.TODO(), tokenString, "blacklisted", time.Until(time.Unix(expirationTime, 0))).Err()
		if err == nil {
			return nil // Operation was successful, return
		}

		currentRetry++
		time.Sleep(2 * time.Second) // Wait before retrying
	}

	if currentRetry == maxRetries {
		logrus.Printf("Max retries reached, could not blacklist the token")
		return fmt.Errorf("max retries reached, could not blacklist the token")
	}

	return nil // This line is technically unreachable but added for completeness
}

// CheckRateLimit checks the rate limit for a given IP in Redis
func CheckRateLimit(ip string, rdb *redis.Client) (bool, error) {
	ctx := context.TODO()
	// Try to increment the count for this IP
	newCount, err := rdb.Incr(ctx, ip).Result()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"ip": ip,
		}).Errorf("Redis error: %v", err)
		return false, err
	}

	// If this is the first request from this IP, set the key to expire after 1 minute
	if newCount == 1 {
		if _, err := rdb.Expire(ctx, ip, 1*time.Minute).Result(); err != nil {
			logrus.WithFields(logrus.Fields{
				"ip": ip,
			}).Errorf("Redis error: %v", err)
			return false, err
		}
	}

	// Check if the IP has exceeded the limit (e.g., 10 requests per minute)
	if newCount > 10 {
		logrus.WithFields(logrus.Fields{
			"ip": ip,
		}).Warn("Rate limit exceeded")
		return false, nil
	}

	return true, nil
}
