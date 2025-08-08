// Package postgres defines error codes for PostgreSQL database operations.
// Check: https://www.postgresql.org/docs/18/errcodes-appendix.html
//
// This package is inspired by: https://github.com/michaljemala/pqerror
package postgres

import "github.com/lib/pq"

type Error string

// IsErrCode checks if the error is a PostgreSQL error with the specified code.
func IsErrCode(err error, code pq.ErrorCode) bool {
	if err == nil {
		return false
	}

	pqErr, ok := err.(*pq.Error)
	if !ok {
		return false
	}

	return pqErr.Code == code
}

const (
	// Class 23 — Integrity Constraint Violation
	ErrIntegrityConstraintViolation pq.ErrorCode = "23000"
	ErrNotNullViolation             pq.ErrorCode = "23502"
	ErrForeignKeyViolation          pq.ErrorCode = "23503"
	ErrUniqueViolation              pq.ErrorCode = "23505"
	ErrCheckViolation               pq.ErrorCode = "23514"
)
