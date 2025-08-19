package cron

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type postgresLocker struct {
	db            *sql.DB
	lockTableName string
	lockTimeout   time.Duration
}

func NewPostgresLocker(db *sql.DB) (gocron.Locker, error) {
	locker := &postgresLocker{
		db:            db,
		lockTableName: "cron_locks",
		lockTimeout:   30 * time.Second,
	}

	return locker, nil
}

func (l *postgresLocker) Lock(ctx context.Context, key string) (gocron.Lock, error) {
	lock := &postgresLock{
		db:        l.db,
		key:       key,
		tableName: l.lockTableName,
		timeout:   l.lockTimeout,
		lockedBy:  generateLockID(),
	}

	if err := lock.acquire(ctx); err != nil {
		return nil, err
	}

	return lock, nil
}

type postgresLock struct {
	db        *sql.DB
	key       string
	tableName string
	timeout   time.Duration
	lockedBy  string
}

func (l *postgresLock) acquire(ctx context.Context) error {
	expiresAt := time.Now().UTC().Add(l.timeout)

	// First, try to insert a new lock or update an expired one
	query := fmt.Sprintf(`
		INSERT INTO %s (key, locked_by, locked_at, expires_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP, $3)
		ON CONFLICT (key) DO UPDATE
		SET locked_by = $2, locked_at = CURRENT_TIMESTAMP, expires_at = $3
		WHERE %s.expires_at < CURRENT_TIMESTAMP
		RETURNING key
	`, l.tableName, l.tableName)

	var returnedKey string
	err := l.db.QueryRowContext(ctx, query, l.key, l.lockedBy, expiresAt).Scan(&returnedKey)

	if err == sql.ErrNoRows {
		// Lock exists and is not expired
		return fmt.Errorf("failed to acquire lock: lock is held by another process")
	} else if err != nil {
		return fmt.Errorf("failed to acquire lock: %w", err)
	}

	return nil
}

func (l *postgresLock) Unlock(ctx context.Context) error {
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE key = $1 AND locked_by = $2
	`, l.tableName)

	result, err := l.db.ExecContext(ctx, query, l.key, l.lockedBy)
	if err != nil {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("lock not found or not owned by this instance")
	}

	return nil
}

func generateLockID() string {
	return fmt.Sprintf("%s-%d", getHostname(), time.Now().UnixNano())
}

func getHostname() string {
	hostname := "unknown"
	if h, err := os.Hostname(); err == nil {
		hostname = h
	}
	return hostname
}
