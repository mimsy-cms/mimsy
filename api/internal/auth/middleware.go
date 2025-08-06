package auth

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	auth_interface "github.com/mimsy-cms/mimsy/internal/interfaces/auth"
)

type contextKey string

const UserContextKey contextKey = "user"

func WithUser(authService auth_interface.AuthService) func(http.Handler) http.Handler {
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

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) *auth_interface.User {
	if user, ok := ctx.Value(UserContextKey).(*auth_interface.User); ok {
		return user
	}
	return nil
}

type DBWrapper struct {
	DB *sql.DB
}

func (w *DBWrapper) QueryRow(query string, args ...any) auth_interface.Row {
	return &RowWrapper{w.DB.QueryRow(query, args...)}
}

func (w *DBWrapper) QueryRowContext(ctx context.Context, query string, args ...any) auth_interface.Row {
	return &RowWrapper{w.DB.QueryRowContext(ctx, query, args...)}
}

func (w *DBWrapper) Exec(query string, args ...interface{}) (sql.Result, error) {
	return w.DB.Exec(query, args...)
}

func (w *DBWrapper) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return w.DB.ExecContext(ctx, query, args...)
}

func (w *DBWrapper) Query(query string, args ...any) (*sql.Rows, error) {
	return w.DB.Query(query, args...)
}

func (w *DBWrapper) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return w.DB.QueryContext(ctx, query, args...)
}

type RowWrapper struct {
	Row *sql.Row
}

func (r *RowWrapper) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}
