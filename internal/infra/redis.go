// internal/infra/redis.go
package infra

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient holds the active connection pool.
var RedisClient *redis.Client

// InitRedis establishes the connection and pings the server to ensure it is alive.
func InitRedis(addr string, username string, password string) {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     username,
		Password:     password,
		DB:           0,
		MinIdleConns: 5,
	})

	// Ping the database to verify connectivity on startup
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("[FATAL] Failed to connect to Redis at %s: %v", addr, err)
	}

	log.Println("[INFO] Successfully connected to Redis")
}
