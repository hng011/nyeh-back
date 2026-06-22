package auth

import (
	"net/http"
	"nyeh-back/internal/core"
	da "nyeh-back/internal/domain/auth"
	"time"
)

type AuthHandler struct {
	sessionCache da.UserSession
}

func NewAuthHandler(sessionCache da.UserSession) *AuthHandler {
	return &AuthHandler{
		sessionCache: sessionCache,
	}
}

// setAccessTokenCookie to standardize access_token cookie creation
func setAccessTokenCookie(w http.ResponseWriter, r *http.Request, token string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     core.ACCESS_TOKEN_COOKIE_KEY,
		Value:    token,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		HttpOnly: true,
		Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteStrictMode,
	})
}

// setRefreshTokenCookie to standardize refresh_cookie cookie creation
func setRefreshTokenCookie(w http.ResponseWriter, r *http.Request, token string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     core.REFRESH_TOKEN_COOKIE_KEY,
		Value:    token,
		Path:     "/",
		MaxAge:   int(ttl.Seconds()),
		HttpOnly: true,
		Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteStrictMode,
	})
}
