package auth

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"nyeh-back/internal/core"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Helper to construct the OAuth configuration dynamically per request
func getOAuthConfig(r *http.Request) *oauth2.Config {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	return &oauth2.Config{
		ClientID:     core.Settings.GOOGLE_OAUTH_CLIENT_ID,
		ClientSecret: core.Settings.GOOGLE_OAUTH_CLIENT_SECRET,
		RedirectURL:  fmt.Sprintf("%s/api/v1/auth/google/callback", baseURL),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

// GoogleLoginHandler godoc
//
//	@Summary	Handle Google Login with SSO
//	@Accept		json
//	@Produce	json
//	@Router		/auth/google/login [get]
func (h *AuthHandler) GoogleLoginHandler(w http.ResponseWriter, r *http.Request) {
	oauthConf := getOAuthConfig(r)

	fmt.Println("DEBUG - Sending this Redirect URI to Google:", oauthConf.RedirectURL)

	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	// Production tracking: Save state to a secure, short-lived HTTP-only cookie to validate during callback
	cookie := &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(3 * time.Minute),
		HttpOnly: true,
		Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	authURL := oauthConf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// GoogleCallbackHandler godoc
//
//	@Summary	Handle Google Callback for SSO Access Token
//	@Accept		json
//	@Produce	json
//	@Router		/auth/google/callback [get]
func (h *AuthHandler) GoogleCallbackHandler(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || r.FormValue("state") != stateCookie.Value {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", MaxAge: -1, Path: "/"}) // Remove GOOGLE OAUTH Cookie

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "Code not found", http.StatusBadRequest)
		return
	}

	oauthConf := getOAuthConfig(r)
	token, err := oauthConf.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, fmt.Sprintf("Token exchange failed: %v", err), http.StatusInternalServerError)
		return
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		http.Error(w, "Failed to fetch user data from Google", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read user data response", http.StatusInternalServerError)
		return
	}

	var googleUser struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(data, &googleUser); err != nil {
		http.Error(w, "Failed to parse user data", http.StatusInternalServerError)
		return
	}

	if googleUser.Email != core.Settings.GOOGLE_ALLOWED_EMAIL {
		http.Error(w, "Unauthorized: Unable to Login", http.StatusForbidden)
		return
	}

	// 6. Generate the Session Access Token (JWT)
	core.GenerateAccessToken(w, googleUser.Email, r.Host)

	// 7. Generate refresh token
	TTL_REFRESH_TOKEN := core.Settings.TTL_REFRESH_TOKEN_HOURS * uint(time.Hour)

	refreshToken, err := core.GenerateRefreshToken(w, time.Duration(TTL_REFRESH_TOKEN))
	if err != nil {
		http.Error(w, "Failed to generate refresh token", http.StatusInternalServerError)
		return
	}

	// 8. Save refresh token to redis
	err = h.sessionCache.SetSession(r.Context(), refreshToken, googleUser.Email, time.Duration(TTL_REFRESH_TOKEN))
	if err != nil {
		http.Error(w, "Failed to save session", http.StatusInternalServerError)
		return
	}

	// 9. Return the JWT token to the client
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "authenticated",
	})
}
