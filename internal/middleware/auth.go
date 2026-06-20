package middleware

import (
	"context"
	"net/http"
	c "nyeh-back/internal/core"
	"strings"
)

// contextKey is a custom type to prevent context key collisions
type contextKey string

const UserEmailKey contextKey = "user_email"

// RequireAuth is the middleware that protects secure routes
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Extract the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// 2. Ensure it follows the "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Unauthorized: Invalid Authorization format", http.StatusUnauthorized)
			return
		}
		tokenString := parts[1]

		// 3. Validate the JWT using the auth package we built earlier
		claims, err := c.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Unauthorized: Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// 4. Inject the validated email into the request context
		ctx := context.WithValue(r.Context(), UserEmailKey, claims.Subject)

		// 5. Create a new request with the updated context and pass it to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
