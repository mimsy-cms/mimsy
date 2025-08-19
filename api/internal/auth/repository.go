package auth

import (
	"context"
	"time"

	"github.com/mimsy-cms/mimsy/internal/config"
)

type User struct {
	ID                 int64     `json:"id"`
	Email              string    `json:"email"`
	PasswordHash       string    `json:"-"`
	IsAdmin            bool      `json:"is_admin"`
	MustChangePassword bool      `json:"must_change_password"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type Repository interface {
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
	GetUsers(ctx context.Context) ([]User, error)
}

type repositoy struct{}

func NewRepository() Repository {
	return &repositoy{}
}

func (r *repositoy) CountUsers(ctx context.Context) (int, error) {
	var count int
	err := config.GetDB(ctx).QueryRowContext(ctx, `SELECT COUNT(*) FROM "user"`).Scan(&count)
	return count, err
}

func (r *repositoy) InsertUser(ctx context.Context, email, password string, isAdmin, mustChange bool) error {
	_, err := config.GetDB(ctx).ExecContext(ctx,
		`INSERT INTO "user" (email, password, is_admin, must_change_password, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6)`,
		email, password, isAdmin, mustChange, time.Now(), time.Now())
	return err
}

func (r *repositoy) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var u User
	err := config.GetDB(ctx).QueryRowContext(ctx, `SELECT id, email, password, is_admin, must_change_password, created_at, updated_at FROM "user" WHERE email = $1`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.MustChangePassword, &u.CreatedAt, &u.UpdatedAt)
	return &u, err
}

func (r *repositoy) DeleteExpiredSessions(ctx context.Context) error {
	_, err := config.GetDB(ctx).ExecContext(ctx, `DELETE FROM session WHERE expires_at < NOW()`)
	return err
}

func (r *repositoy) CreateSession(ctx context.Context, token string, userID int64, expiresAt time.Time) error {
	_, err := config.GetDB(ctx).ExecContext(ctx,
		`INSERT INTO session (id, user_id, expires_at) VALUES ($1, $2, $3)`,
		token, userID, expiresAt)
	return err
}

func (r *repositoy) DeleteSession(ctx context.Context, token string) error {
	_, err := config.GetDB(ctx).ExecContext(ctx, `DELETE FROM session WHERE id = $1`, token)
	return err
}

func (r *repositoy) GetUserPassword(ctx context.Context, userID int64) (string, error) {
	var hash string
	err := config.GetDB(ctx).QueryRowContext(ctx, `SELECT password FROM "user" WHERE id = $1`, userID).Scan(&hash)
	return hash, err
}

func (r *repositoy) UpdatePassword(ctx context.Context, userID int64, newHash string) error {
	_, err := config.GetDB(ctx).ExecContext(ctx, `UPDATE "user" SET password = $1, must_change_password = FALSE, updated_at = $2 WHERE id = $3`, newHash, time.Now(), userID)
	return err
}

func (r *repositoy) UserExists(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := config.GetDB(ctx).QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM "user" WHERE email=$1)`, email).Scan(&exists)
	return exists, err
}

func (r *repositoy) GetUserBySessionToken(ctx context.Context, token string) (*User, error) {
	var user User
	var userID int

	err := config.GetDB(ctx).QueryRowContext(ctx, `SELECT user_id FROM session WHERE id = $1`, token).Scan(&userID)
	if err != nil {
		return nil, err
	}

	err = config.GetDB(ctx).QueryRowContext(ctx, `
		SELECT id, email, password, must_change_password, is_admin, created_at, updated_at
		FROM "user"
		WHERE id = $1`, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.MustChangePassword, &user.IsAdmin, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repositoy) GetUsers(ctx context.Context) ([]User, error) {
	rows, err := config.GetDB(ctx).QueryContext(ctx, `SELECT id, email, password, is_admin, must_change_password, created_at, updated_at FROM "user"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.IsAdmin, &u.MustChangePassword, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	return users, rows.Err()
}
