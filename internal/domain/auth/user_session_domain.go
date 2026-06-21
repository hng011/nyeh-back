package auth

import (
	"context"
	"time"
)

type UserSession interface {
	SetSession(ctx context.Context, refreshToken string, email string, expiresIn time.Duration) error
	GetSession(ctx context.Context, refreshToken string) (string, time.Duration, error)
	DeleteSession(ctx context.Context, refreshToken string) error
}
