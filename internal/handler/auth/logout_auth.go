package auth

import (
	"encoding/json"
	"net/http"
	"nyeh-back/internal/core"
)

// LogoutHandler godoc
//
//	@Summary	Logout handler
//	@Accept		json
//	@Produce	json
//	@Router		/auth/logout [post]
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := r.Cookie(core.REFRESH_TOKEN_COOKIE_KEY)
	if err == nil && refreshToken.Value != "" {
		// Only hits Redis if the cookie existed in the request
		_ = h.sessionCache.DeleteSession(r.Context(), core.HashToken(refreshToken.Value))
	}

	// Delete refresh_token from cookie
	setRefreshTokenCookie(w, r, "", -1)

	// Delete access_token from cookie
	setAccessTokenCookie(w, r, "", -1)

	// 4. Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Successfully Logged Out!",
	}
	json.NewEncoder(w).Encode(response)
}
