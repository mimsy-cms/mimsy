package auth

import (
	"context"
	"database/sql"
	"net/http"
	"strings"
)

type contextKey string

const userContextKey contextKey = "user"

func WithUser(db *sql.DB) func(http.Handler) http.Handler {
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

			var userID int
			err := db.QueryRow(`SELECT user_id FROM "session" WHERE id = $1`, token).Scan(&userID)
			if err != nil {
				if err != sql.ErrNoRows {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				next.ServeHTTP(w, r)
				return
			}

			user := User{}

			err = db.QueryRow(`SELECT id, email, password, must_change_password, is_admin FROM "user" WHERE id = $1`, userID).Scan(&user.ID, &user.Email, &user.PasswordHash, &user.MustChangePassword, &user.IsAdmin)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, &user)
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
