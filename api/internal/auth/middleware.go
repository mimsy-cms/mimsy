package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey contextKey = "user"

func WithUser(authService AuthService) func(http.Handler) http.Handler {
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

			if token == "" {
				next.ServeHTTP(w, r)
				return
			}

			user, err := authService.GetUserBySessionToken(r.Context(), token)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) *User {
	if user, ok := ctx.Value(userContextKey).(*User); ok {
		return user
	}
	return nil
}
