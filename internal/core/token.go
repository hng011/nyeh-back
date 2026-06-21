package core

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const REFRESH_TOKEN_COOKIE_KEY string = "refresh_token"
const ACCESS_TOKEN_COOKIE_KEY string = "access_token"

// GenerateToken creates a signed HS256 JWT for a whitelisted email
func GenerateAccessToken(w http.ResponseWriter, email string, host string) {
	secretKey := []byte(Settings.JWT_AUTH_TOKEN)
	if len(secretKey) == 0 {
		http.Error(w, "Secret Key is Invalid", http.StatusInternalServerError)
		return
	}

	TTL := time.Duration(Settings.TTL_ACCESS_TOKEN_MINUTES) * time.Minute
	var claims jwt.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TTL)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    host,
		Subject:   email,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		http.Error(w, "Failed to generate session token", http.StatusInternalServerError)
		return
	}

	// STORE ACCESS_TOKEN TO COOKIE
	http.SetCookie(w, &http.Cookie{
		Name:     ACCESS_TOKEN_COOKIE_KEY,
		Value:    token,
		Path:     "/",
		Expires:  time.Now().Add(time.Duration(TTL)),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

// ValidateToken parses and verifies the signature and expiration of a token string
func ValidateToken(w http.ResponseWriter, r *http.Request) (*jwt.RegisteredClaims, error) {
	secretKey := []byte(Settings.JWT_AUTH_TOKEN)

	// 1. Get Refresh Token from httpOnly cookie
	tokenStr, err := r.Cookie(ACCESS_TOKEN_COOKIE_KEY)
	if err != nil {
		http.Error(w, "Unauthorized: No access token", http.StatusUnauthorized)
		return nil, err
	}

	token, err := jwt.ParseWithClaims(tokenStr.Value, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
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

// GenerateRefreshToken creates a random 32-byte opaque token
func GenerateRefreshToken(w http.ResponseWriter, ttl time.Duration) (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	refresh_token := base64.URLEncoding.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:     REFRESH_TOKEN_COOKIE_KEY,
		Value:    refresh_token,
		Path:     "/",
		Expires:  time.Now().Add(ttl),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	return refresh_token, nil
}
