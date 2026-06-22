package auth

import (
	"encoding/json"
	"net/http"
	"nyeh-back/internal/core"
)

// RefreshHandler godoc
//
//	@Summary	Refresh token rotation handler
//	@Accept		json
//	@Produce	json
//	@Router		/auth/refresh [post]
func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	// Get Refresh Token from httpOnly cookie
	cookie, err := r.Cookie(core.REFRESH_TOKEN_COOKIE_KEY)
	if err != nil {
		http.Error(w, "Unauthorized: No refresh_token", http.StatusUnauthorized)
		return
	}
	oldRefreshToken := core.HashToken(cookie.Value)

	// validate against Redis
	email, remainingTTL, err := h.sessionCache.GetSession(r.Context(), oldRefreshToken)
	if err != nil {
		http.Error(w, "Unauthorized: No refresh_token", http.StatusUnauthorized)
		return
	}

	if email != core.Settings.GOOGLE_ALLOWED_EMAIL {
		http.Error(w, "Unauthorized refresh_token", http.StatusUnauthorized)
		return
	}

	// Security: Rotate the Refresh Token
	// Delete the old one and create a new one
	_ = h.sessionCache.DeleteSession(r.Context(), oldRefreshToken)
	newRefreshToken, newHashedRefreshToken, _, err := core.GenerateRefreshToken()
	if err != nil {
		http.Error(w, "Failed to generate a refresh_token", http.StatusUnauthorized)
		return
	}

	// Store the session back to cache storage
	err = h.sessionCache.SetSession(r.Context(), newHashedRefreshToken, email, remainingTTL)
	if err != nil {
		http.Error(w, "Failed to store a refresh_token to cache", http.StatusInternalServerError)
		return
	}

	// Store new refresh token to cookie
	setRefreshTokenCookie(w, r, newRefreshToken, remainingTTL)

	// Generate new Access Token
	accessToken, ttlAccessToken, err := core.GenerateAccessToken(email, r.Host)
	if err != nil {
		http.Error(w, "Failed to generate an access_token", http.StatusUnauthorized)
		return
	}
	// Store Access Token to cookie
	setAccessTokenCookie(w, r, accessToken, ttlAccessToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "access token generated",
	})

}
