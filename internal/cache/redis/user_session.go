package redis

import (
	"context"
	"errors"
	"fmt"
	da "nyeh-back/internal/domain/auth"
	"time"

	"github.com/redis/go-redis/v9"
)

const session_key_name string = "refresh_token"

type userSessionCache struct {
	client *redis.Client
}

func NewUserSessionCache(client *redis.Client) da.UserSession {
	return &userSessionCache{
		client: client,
	}
}

func (c *userSessionCache) SetSession(ctx context.Context, refreshToken string, email string, expiresIn time.Duration) error {
	// Prefix the key to organize redis
	return c.client.Set(ctx, fmt.Sprintf("%v:%v", session_key_name, refreshToken), email, expiresIn).Err()
}

func (c *userSessionCache) GetSession(ctx context.Context, refreshToken string) (string, time.Duration, error) {
	key := fmt.Sprintf("%v:%v", session_key_name, refreshToken)

	// 1. Get Session
	val, err := c.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", 0, errors.New("Refresh token not found")
	} else if err != nil {
		return "", 0, err
	}

	// 2. Get the exact remaining TTL of this specific key
	ttl, err := c.client.TTL(ctx, key).Result()
	if err != nil {
		return "", 0, err
	}

	// If TTL is negative, it means the key exists but has no expiration (shouldn't happen here)
	// or it was deleted microseconds after our Get command.
	if ttl <= 0 {
		return "", 0, errors.New("refresh token expired")
	}

	return val, ttl, err
}

func (c *userSessionCache) DeleteSession(ctx context.Context, refreshToken string) error {
	return c.client.Del(ctx, fmt.Sprintf("%v:%v", session_key_name, refreshToken)).Err()
}
