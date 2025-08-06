package auth_interface

import (
	"context"
	"database/sql"
)

type DB interface {
	QueryRow(query string, args ...any) Row
	QueryRowContext(ctx context.Context, query string, args ...any) Row

	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)

	Query(query string, args ...any) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

type Row interface {
	Scan(dest ...interface{}) error
}

type AuthService interface {
	GetUserBySessionToken(ctx context.Context, token string) (*User, error)
}

type User struct {
	ID                 int64  `json:"id"`
	Email              string `json:"email"`
	PasswordHash       string `json:"-"`
	IsAdmin            bool   `json:"is_admin"`
	MustChangePassword bool   `json:"must_change_password"`
}
