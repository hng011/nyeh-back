package middleware

import (
	"context"
	"net/http"
	"nyeh-back/internal/core"
)

// contextKey is a custom type to prevent context key collisions

// RequireAuth is the middleware that protects secure routes
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// 1. Get Refresh Token from httpOnly cookie
		tokenStr, err := r.Cookie(core.ACCESS_TOKEN_COOKIE_KEY)
		if err != nil {
			http.Error(w, "Unauthorized: No access token", http.StatusUnauthorized)
			return
		}

		// 2. Validate the cookie
		claims, err := core.ValidateToken(tokenStr.Value)
		if err != nil {
			http.Error(w, "Unauthorized: No access token", http.StatusUnauthorized)
			return
		}

		// 3. Validate the email
		var email string = claims.Subject
		if email != core.Settings.ALLOWED_EMAIL {
			http.Error(w, "Unauthorized access_token", http.StatusUnauthorized)
			return
		}

		// 3. Inject the identity
		ctx := context.WithValue(r.Context(), core.JWT_CLAIM_USER_EMAIL_KEY, email)

		// 4. Continue
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
