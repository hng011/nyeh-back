package handler_auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"nyeh-back/internal/core"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleSSOHandler struct{}

func (sso *GoogleSSOHandler) GetOAuthConfig(scheme string, host string) *oauth2.Config {

	baseURL := fmt.Sprintf("%s://%s", scheme, host)

	return &oauth2.Config{
		ClientID:     core.Settings.GOOGLE_OAUTH_CLIENT_ID,
		ClientSecret: core.Settings.GOOGLE_OAUTH_CLIENT_SECRET,
		RedirectURL:  fmt.Sprintf("%s/api/v1/auth/google/callback", baseURL),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func (sso *GoogleSSOHandler) FetchEmail(ctx context.Context, accessToken string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://www.googleapis.com/oauth2/v2/userinfo?access_token="+accessToken, nil)
	if err != nil {
		return "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var googleUser struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(data, &googleUser); err != nil {
		return "", err
	}
	return googleUser.Email, nil
}
