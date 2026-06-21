package auth

import (
	"encoding/json"
	"net/http"
	"nyeh-back/internal/core"
	da "nyeh-back/internal/domain/auth"
)

type AuthHandler struct {
	sessionCache da.UserSession
}

func NewAuthHandler(sessionCache da.UserSession) *AuthHandler {
	return &AuthHandler{
		sessionCache: sessionCache,
	}
}

// RefreshHandler godoc
//
//	@Summary	Refresh token rotation handler
//	@Accept		json
//	@Produce	json
//	@Router		/auth/refresh [post]
func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Get Refresh Token from httpOnly cookie
	cookie, err := r.Cookie(core.REFRESH_TOKEN_COOKIE_KEY)
	if err != nil {
		http.Error(w, "Unauthorized: No refresh token", http.StatusUnauthorized)
		return
	}
	oldRefreshToken := core.HashToken(cookie.Value)

	// 2. validate against Redis
	email, remainingTTL, err := h.sessionCache.GetSession(r.Context(), oldRefreshToken)
	if err != nil {
		http.Error(w, "Unauthorized: No refresh token", http.StatusUnauthorized)
		return
	}

	// 3. Security: Rotate the Refresh Token
	// Delete the old one and create a new one to prevent replay attacks
	_ = h.sessionCache.DeleteSession(r.Context(), oldRefreshToken)
	NewRefreshToken, _ := core.GenerateRefreshToken(w, remainingTTL)
	_ = h.sessionCache.SetSession(r.Context(), NewRefreshToken, email, remainingTTL)

	// 4. Generate new Access Token
	core.GenerateAccessToken(w, email, r.Host)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "access token generated",
	})

}
