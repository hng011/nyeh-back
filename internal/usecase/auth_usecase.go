package usecase

import (
	"context"
	"errors"
	"nyeh-back/internal/core"
	da "nyeh-back/internal/domain/auth"
	"time"
)

type authUsecase struct {
	userSessionRepo da.UserSessionDomain
}

func NewAuthUsecase(userSessionRepo da.UserSessionDomain) da.AuthDomain {
	return &authUsecase{
		userSessionRepo: userSessionRepo,
	}
}

func (a *authUsecase) Login(ctx context.Context, email string, host string) (string, time.Duration, string, time.Duration, error) {
	if email != core.Settings.ALLOWED_EMAIL {
		return "", 0, "", 0, errors.New("Unauthorized email")
	}

	accToken, accTTL, err := core.GenerateAccessToken(email, host)
	if err != nil {
		return "", 0, "", 0, err
	}

	rawRefresh, hashedRefresh, refTTL, err := core.GenerateRefreshToken()
	if err != nil {
		return "", 0, "", 0, err
	}

	err = a.userSessionRepo.SetSession(ctx, hashedRefresh, email, refTTL)
	if err != nil {
		return "", 0, "", 0, err
	}

	return accToken, accTTL, rawRefresh, refTTL, nil
}

func (a *authUsecase) Refresh(ctx context.Context, oldHashedRefreshToken string, host string) (string, time.Duration, string, time.Duration, error) {
	email, remainingTTL, err := a.userSessionRepo.GetSession(ctx, oldHashedRefreshToken)
	if err != nil {
		return "", 0, "", 0, err
	}

	if email != core.Settings.ALLOWED_EMAIL {
		return "", 0, "", 0, errors.New("Unauthorized email")
	}

	// Invalidate the old token immediately (Token Rotation)
	_ = a.userSessionRepo.DeleteSession(ctx, oldHashedRefreshToken)

	accToken, accTTL, err := core.GenerateAccessToken(email, host)
	if err != nil {
		return "", 0, "", 0, err
	}

	newRawRefresh, newHashedRefresh, _, err := core.GenerateRefreshToken()
	if err != nil {
		return "", 0, "", 0, err
	}

	err = a.userSessionRepo.SetSession(ctx, newHashedRefresh, email, remainingTTL)
	if err != nil {
		return "", 0, "", 0, err
	}

	return accToken, accTTL, newRawRefresh, remainingTTL, nil
}

func (a *authUsecase) Logout(ctx context.Context, oldRawRefreshToken string) {
	_ = a.userSessionRepo.DeleteSession(ctx, core.HashToken(oldRawRefreshToken))
}
