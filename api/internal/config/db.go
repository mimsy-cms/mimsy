package config

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
)

// DB abstracts the database operations to allow the use of either sql.DB or sql.Tx.
// Comments are taken from the sql package methods.
type DB interface {
	// Exec executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	//
	// Exec uses [context.Background] internally; to specify the context, use
	// [DB.ExecContext].
	Exec(query string, args ...any) (sql.Result, error)
	// ExecContext executes a query without returning any rows.
	// The args are for any placeholder parameters in the query.
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	// Prepare creates a prepared statement for use within a transaction.
	//
	// The returned statement operates within the transaction and will be closed
	// when the transaction has been committed or rolled back.
	//
	// To use an existing prepared statement on this transaction, see [Tx.Stmt].
	//
	// Prepare uses [context.Background] internally; to specify the context, use
	// [Tx.PrepareContext].
	Prepare(query string) (*sql.Stmt, error)
	// PrepareContext creates a prepared statement for use within a transaction.
	//
	// The returned statement operates within the transaction and will be closed
	// when the transaction has been committed or rolled back.
	//
	// To use an existing prepared statement on this transaction, see [Tx.Stmt].
	//
	// The provided context will be used for the preparation of the context, not
	// for the execution of the returned statement. The returned statement
	// will run in the transaction context.
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	// Query executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	//
	// Query uses [context.Background] internally; to specify the context, use
	// [DB.QueryContext].
	Query(query string, args ...any) (*sql.Rows, error)
	// QueryContext executes a query that returns rows, typically a SELECT.
	// The args are for any placeholder parameters in the query.
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	// QueryRow executes a query that is expected to return at most one row.
	// QueryRow always returns a non-nil value. Errors are deferred until
	// [Row]'s Scan method is called.
	// If the query selects no rows, the [*Row.Scan] will return [ErrNoRows].
	// Otherwise, [*Row.Scan] scans the first selected row and discards
	// the rest.
	//
	// QueryRow uses [context.Background] internally; to specify the context, use
	// [DB.QueryRowContext].
	QueryRow(query string, args ...any) *sql.Row
	// QueryRowContext executes a query that is expected to return at most one row.
	// QueryRowContext always returns a non-nil value. Errors are deferred until
	// [Row]'s Scan method is called.
	// If the query selects no rows, the [*Row.Scan] will return [ErrNoRows].
	// Otherwise, [*Row.Scan] scans the first selected row and discards
	// the rest.
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

type contextKey struct{}

// GetDB retrieves the database connection from the context.
func GetDB(ctx context.Context) DB {
	return ctx.Value(contextKey{}).(DB)
}

// getTx retrieves the transaction from the context if it exists.
// This function is not exported as we don't want to expose the transaction
// to the outside world.
func getTx(ctx context.Context) *sql.Tx {
	if tx, ok := ctx.Value(contextKey{}).(*sql.Tx); ok {
		return tx
	}
	return nil
}

// ContextWithDB creates a new context with the provided database connection.
func ContextWithDB(ctx context.Context, db DB) context.Context {
	return context.WithValue(ctx, contextKey{}, db)
}

// WithDB is a middleware that injects the database connection into the request context.
func WithDB(db *sql.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), contextKey{}, db)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// WithinTx executes a function within a transaction context.
// If the context already contains a transaction, it will use that.
// If not, it will create a new transaction and commit it after the function execution.
// If an error is returned inside the fn function, the transaction will be rolled back.
func WithinTx(ctx context.Context, fn func(context.Context) error) error {
	tx := getTx(ctx)
	if tx != nil {
		return fn(ctx)
	}

	db, ok := ctx.Value(contextKey{}).(*sql.DB)
	if !ok {
		return fmt.Errorf("context does not contain a database connection")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	txCtx := context.WithValue(ctx, contextKey{}, tx)

	if err = fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return fmt.Errorf("error during rollback: %w", rollbackErr)
		}
		return err
	}

	return tx.Commit()
}
