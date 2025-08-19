package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserContextKey contextKey = "user"

func WithRequestUser(authService Service) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var token string

			authHeader := r.Header.Get("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				cookie, err := r.Cookie("session")
				if err == nil {
					token = cookie.Value
				}
			}

			user, err := authService.GetUserBySessionToken(r.Context(), token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequestUser retrieves the authenticated user from the request context.
// Returns nil if no user is authenticated.
func RequestUser(ctx context.Context) *User {
	if user, ok := ctx.Value(UserContextKey).(*User); ok {
		return user
	}
	return nil
}
