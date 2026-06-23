package auth

import (
	"context"
	"time"
)

type AuthDomain interface {
	Login(ctx context.Context, email string, host string) (string, time.Duration, string, time.Duration, error)
	Refresh(ctx context.Context, oldRawRefreshToken string, host string) (string, time.Duration, string, time.Duration, error)
	Logout(ctx context.Context, oldRawRefreshToken string)
}
