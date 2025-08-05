package auth_interface

import (
	"context"
	"database/sql"
)

type DB interface {
	QueryRow(query string, args ...any) Row
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

type Row interface {
	Scan(dest ...interface{}) error
}
