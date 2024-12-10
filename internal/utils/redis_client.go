package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

// RedisClient struct holds the Redis client and methods
type RedisClient struct {
	Client *redis.Client
}

// NewRedisClient initializes and returns a RedisClient
func NewRedisClient() *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Test the connection
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis is connected successfully")
	return &RedisClient{Client: client}
}

// Set stores a value in Redis
func (r *RedisClient) Set(key, value string, expiration time.Duration) error {
	ctx := context.Background()
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// Get retrieves a value from Redis
func (r *RedisClient) Get(key string) (string, error) {
	ctx := context.Background()
	val, err := r.Client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// Delete deletes a key from Redis
func (r *RedisClient) Delete(key string) error {
	ctx := context.Background()
	return r.Client.Del(ctx, key).Err()
}
