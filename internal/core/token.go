package core

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a signed HS256 JWT for a whitelisted email
func GenerateAccessToken(email string, host string) (string, time.Duration, error) {
	secretKey := []byte(Settings.JWT_AUTH_TOKEN)
	if len(secretKey) == 0 {
		return "", 0, errors.New("Unable to find the JWT_AUTH_TOKEN")
	}

	TTL := time.Duration(Settings.TTL_ACCESS_TOKEN_MINUTES) * time.Minute
	var claims jwt.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    host,
		Subject:   email,
	}

	access_token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		return "", 0, errors.New("Failed to generate session token")
	}

	return access_token, TTL, nil
}

// ValidateToken parses and verifies the signature and expiration of a token string
func ValidateToken(tokenStr string) (*jwt.RegisteredClaims, error) {

	// Validate the JWT token with the secret key
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
		// Validate the signing algorithm is HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(Settings.JWT_AUTH_TOKEN), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// GenerateRefreshToken creates a random 32-byte opaque token
func GenerateRefreshToken() (string, string, time.Duration, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", "", 0, err
	}

	ttlRefreshToken := time.Duration(Settings.TTL_REFRESH_TOKEN_HOURS) * time.Hour
	rawRefreshToken := base64.URLEncoding.EncodeToString(b)
	hashedRefreshToken := HashToken(rawRefreshToken)

	return rawRefreshToken, hashedRefreshToken, ttlRefreshToken, nil
}
