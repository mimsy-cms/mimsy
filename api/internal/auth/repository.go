package auth

import (
	"context"
	"database/sql"
	"time"
)

type User struct {
	ID                 int64
	Email              string
	PasswordHash       string
	IsAdmin            bool
	MustChangePassword bool
}

type AuthRepository interface {
	CountUsers(ctx context.Context) (int, error)
	InsertUser(ctx context.Context, email, password string, isAdmin, mustChange bool) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	DeleteExpiredSessions(ctx context.Context) error
	CreateSession(ctx context.Context, token string, userID int64, expiresAt time.Time) error
	DeleteSession(ctx context.Context, token string) error
	GetUserPassword(ctx context.Context, userID int64) (string, error)
	UpdatePassword(ctx context.Context, userID int64, newHash string) error
	UserExists(ctx context.Context, email string) (bool, error)
	GetUserBySessionToken(ctx context.Context, token string) (*User, error)
}

type authRepo struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepo{db}
}

func (r *authRepo) CountUsers(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM "user"`).Scan(&count)
	return count, err
}

func (r *authRepo) InsertUser(ctx context.Context, email, password string, isAdmin, mustChange bool) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO "user" (email, password, is_admin, must_change_password) VALUES ($1, $2, $3, $4)`,
		email, password, isAdmin, mustChange)
	return err
}

func (r *authRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := r.db.QueryRowContext(ctx, `SELECT id, email, password, is_admin, must_change_password FROM "user" WHERE email = $1`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.MustChangePassword)
	return &u, err
}

func (r *authRepo) DeleteExpiredSessions(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM session WHERE expires_at < NOW()`)
	return err
}

func (r *authRepo) CreateSession(ctx context.Context, token string, userID int64, expiresAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO session (id, user_id, expires_at) VALUES ($1, $2, $3)`,
		token, userID, expiresAt)
	return err
}

func (r *authRepo) DeleteSession(ctx context.Context, token string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM session WHERE id = $1`, token)
	return err
}

func (r *authRepo) GetUserPassword(ctx context.Context, userID int64) (string, error) {
	var hash string
	err := r.db.QueryRowContext(ctx, `SELECT password FROM "user" WHERE id = $1`, userID).Scan(&hash)
	return hash, err
}

func (r *authRepo) UpdatePassword(ctx context.Context, userID int64, newHash string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE "user" SET password = $1, must_change_password = FALSE WHERE id = $2`, newHash, userID)
	return err
}

func (r *authRepo) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM "user" WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}

func (r *authRepo) GetUserBySessionToken(ctx context.Context, token string) (*User, error) {
	var user User
	var userID int

	err := r.db.QueryRowContext(ctx, `SELECT user_id FROM session WHERE id = $1`, token).Scan(&userID)
	if err != nil {
		return nil, err
	}

	err = r.db.QueryRowContext(ctx, `
		SELECT id, email, password, must_change_password, is_admin
		FROM "user"
		WHERE id = $1`, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.MustChangePassword, &user.IsAdmin,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
