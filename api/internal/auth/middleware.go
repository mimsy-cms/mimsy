package auth

import (
	"context"
	"net/http"
	"strings"

	auth_interface "github.com/mimsy-cms/mimsy/internal/interfaces/auth"
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

type DBWrapper struct {
	DB *sql.DB
}

func (w *DBWrapper) QueryRow(query string, args ...any) auth_interface.Row {
	return w.DB.QueryRow(query, args...)
}

func (w *DBWrapper) QueryRowContext(ctx context.Context, query string, args ...any) auth_interface.Row {
	return w.DB.QueryRowContext(ctx, query, args...)
}

func (w *DBWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	return w.DB.Exec(query, args...)
}

func (w *DBWrapper) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return w.DB.ExecContext(ctx, query, args...)
}

type RowWrapper struct {
	Row *sql.Row
}

func (r *RowWrapper) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}
