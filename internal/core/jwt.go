package core

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a signed HS256 JWT for a whitelisted email
func GenerateToken(email string, host string) (string, error) {
	secretKey := []byte(Settings.JWT_AUTH_TOKEN)
	if len(secretKey) == 0 {
		return "", errors.New("JWT secret key is not configured")
	}

	var claims jwt.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    host,
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}

// ValidateToken parses and verifies the signature and expiration of a token string
func ValidateToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	secretKey := []byte(Settings.JWT_AUTH_TOKEN)

	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		// Validate the signing algorithm is HS256
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
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
