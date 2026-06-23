package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"nyeh-back/internal/core"
	da "nyeh-back/internal/domain/auth"
	"time"

	"golang.org/x/oauth2"
)

// DTO
type AuthResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type SSOProviderInterface interface {
	GetOAuthConfig(scheme string, host string) *oauth2.Config
	FetchEmail(ctx context.Context, accessToken string) (string, error)
}

type AuthHandler struct {
	authUsecase da.AuthDomain
	providers   map[string]SSOProviderInterface
}

func NewAuthHandler(authUsecase da.AuthDomain) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		providers: map[string]SSOProviderInterface{
			"google": &GoogleSSOHandler{},
		},
	}
}

// SetAccessTokenCookie to standardize access_token cookie creation
func SetAccessTokenCookie(w http.ResponseWriter, r *http.Request, token string, ttl time.Duration) {
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

// SetRefreshTokenCookie to standardize refresh_cookie cookie creation
func SetRefreshTokenCookie(w http.ResponseWriter, r *http.Request, token string, ttl time.Duration) {
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

// LoginHandler godoc
//
//	@Summary	Login handler
//	@Accept		json
//	@Produce	json
//	@Param		provider	path	string	true	"OAuth provider"	Enums(google, github)
//	@Router		/auth/{provider}/login [get]
func (h *AuthHandler) OAuthLoginHandler(w http.ResponseWriter, r *http.Request) {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	provider, exists := h.providers[r.PathValue("provider")]
	if !exists {
		http.Error(w, "Unsupported provider", http.StatusBadRequest)
		return
	}

	oauthConf := provider.GetOAuthConfig(scheme, r.Host)

	b := make([]byte, 32)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)

	http.SetCookie(w, &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   180,
		HttpOnly: true,
		Secure:   r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https",
		SameSite: http.SameSiteLaxMode,
	})
	authURL := oauthConf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// CallbackHandler godoc
//
//	@Summary	Callback handler
//	@Accept		json
//	@Produce	json
//	@Param		provider	path	string	true	"OAuth provider"	Enums(google, github)
//	@Router		/auth/{provider}/callback [get]
func (h *AuthHandler) OAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	stateCookie, err := r.Cookie("oauth_state")
	if err != nil || r.FormValue("state") != stateCookie.Value {
		http.Error(w, "Invalid state parameter", http.StatusBadRequest)
		return
	}
	http.SetCookie(w, &http.Cookie{Name: "oauth_state", MaxAge: -1, Path: "/"})

	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}

	providerName := r.PathValue("provider")
	provider, exists := h.providers[providerName]
	if !exists {
		http.Error(w, "Unsupported provider", http.StatusBadRequest)
		return
	}

	oauthConf := provider.GetOAuthConfig(scheme, r.Host)
	code := r.FormValue("code")
	token, err := oauthConf.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "Token exchange failed", http.StatusInternalServerError)
		return
	}

	email, err := provider.FetchEmail(r.Context(), token.AccessToken)
	if err != nil || email == "" {
		http.Error(w, "Failed to fetch user email", http.StatusInternalServerError)
		return
	}

	accToken, accTTL, refToken, refTTL, err := h.authUsecase.Login(r.Context(), email, r.Host)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	SetAccessTokenCookie(w, r, accToken, accTTL)
	SetRefreshTokenCookie(w, r, refToken, refTTL)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{
		Status:  "success",
		Message: fmt.Sprintf("authenticated via %s", providerName),
	})
}

// RefreshHandler godoc
//
//	@Summary	Refresh Token handler
//	@Accept		json
//	@Produce	json
//	@Router		/auth/refresh [post]
func (h *AuthHandler) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(core.REFRESH_TOKEN_COOKIE_KEY)
	if err != nil {
		http.Error(w, "Unauthorized: Missing refresh token", http.StatusUnauthorized)
		return
	}

	oldHashedToken := core.HashToken(cookie.Value)

	accToken, accTTL, refToken, refTTL, err := h.authUsecase.Refresh(r.Context(), oldHashedToken, r.Host)
	if err != nil {
		http.Error(w, "Unauthorized: Session expired or invalid", http.StatusUnauthorized)
		return
	}

	SetAccessTokenCookie(w, r, accToken, accTTL)
	SetRefreshTokenCookie(w, r, refToken, refTTL)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AuthResponse{Status: "ok", Message: "tokens successfully rotated"})
}

// LogoutHandler godoc
//
//	@Summary	Logout handler
//	@Accept		json
//	@Produce	json
//	@Router		/auth/logout [post]
func (h *AuthHandler) LogoutHandler(w http.ResponseWriter, r *http.Request) {

	refreshToken, err := r.Cookie(core.REFRESH_TOKEN_COOKIE_KEY)
	if err == nil && refreshToken.Value != "" {
		h.authUsecase.Logout(r.Context(), refreshToken.Value)
	}

	SetRefreshTokenCookie(w, r, "", -1)
	SetAccessTokenCookie(w, r, "", -1)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Successfully Logged Out!",
	}
	json.NewEncoder(w).Encode(response)
}
