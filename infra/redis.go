package infra

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client

// InitRedis initializes the Redis client connection.
func InitRedis() error {
	rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	return rdb.Ping(ctx).Err()
}

// GetKey retrieves a value from Redis.
func RdbGetKey(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	return rdb.Get(ctx, key).Result()
}

// SetKey stores a string key-value pair in Redis with a TTL.
func RdbSetKey(key string, value interface{}, ttlSeconds int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()

	return rdb.Set(ctx, key, value, time.Duration(ttlSeconds) * time.Second).Err()
}
