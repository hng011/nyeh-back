package middleware

import (
	"context"
	"net/http"
	"nyeh-back/internal/core"
)

// contextKey is a custom type to prevent context key collisions
type contextKey string

const UserEmailKey contextKey = "user_email"

// RequireAuth is the middleware that protects secure routes
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Validate the cookie
		claims, err := core.ValidateToken(w, r)
		if err != nil {
			// core.ValidateToken already sent the 401 response
			return
		}

		// 2. Inject the identity
		ctx := context.WithValue(r.Context(), UserEmailKey, claims.Subject)

		// 3. Continue
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
