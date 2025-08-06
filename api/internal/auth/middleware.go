package auth

import (
	"context"
	"database/sql"
	"net/http"
	"strings"

	"github.com/mimsy-cms/mimsy/internal/config"
	auth_interface "github.com/mimsy-cms/mimsy/internal/interfaces/auth"
)

type contextKey string

const userContextKey contextKey = "user"

func WithUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var token string

		ctx := r.Context()

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

		var userID int
		err := config.GetDB(ctx).QueryRow(`SELECT user_id FROM "session" WHERE id = $1`, token).Scan(&userID)
		if err != nil {
			if err != sql.ErrNoRows {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			next.ServeHTTP(w, r)
			return
		}

		user := User{}

		err = config.GetDB(ctx).QueryRow(`SELECT id, email, password, must_change_password, is_admin FROM "user" WHERE id = $1`, userID).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.MustChangePassword, &user.IsAdmin)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		userCtx := r.WithContext(context.WithValue(ctx, userContextKey, &user))
		next.ServeHTTP(w, userCtx)
	})
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
