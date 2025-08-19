package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
)

var (
	adminEmail    = os.Getenv("ADMIN_EMAIL")
	adminPassword = os.Getenv("ADMIN_PASSWORD")
)

type LoginResponse struct {
	Session            string `json:"session"`
	MustChangePassword bool   `json:"mustChangePassword"`
}

type Service interface {
	CreateAdminUser(ctx context.Context) error
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	Logout(ctx context.Context, sessionToken string) error
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	Register(ctx context.Context, req CreateUserRequest) error
	GetUserBySessionToken(ctx context.Context, sessionToken string) (*User, error)
	GetUsers(ctx context.Context) ([]User, error)
}

func (s *service) GetUserBySessionToken(ctx context.Context, sessionToken string) (*User, error) {
	return s.authRepository.GetUserBySessionToken(ctx, sessionToken)
}

type service struct {
	authRepository Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) CreateAdminUser(ctx context.Context) error {
	count, err := s.authRepository.CountUsers(ctx)
	if err != nil {
		return fmt.Errorf("failed to count users: %w", err)
	}
	if count > 0 {
		return nil // Admin user already exists
	}

	hash, err := HashPassword(adminPassword)
	if err != nil {
		return fmt.Errorf("failed to hash admin password: %w", err)
	}

	err = s.authRepository.InsertUser(ctx, adminEmail, hash, true, true)
	if err != nil {
		return fmt.Errorf("failed to insert admin user: %w", err)
	}

	return nil
}

func (s *service) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := s.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := CheckPasswordHash(password, user.PasswordHash); err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := s.authRepository.DeleteExpiredSessions(ctx); err != nil {
		return nil, err
	}
	token, err := GenerateSessionToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.authRepository.CreateSession(ctx, token, user.ID, expiresAt); err != nil {
		return nil, err
	}
	return &LoginResponse{Session: token, MustChangePassword: user.MustChangePassword}, nil
}

func (s *service) Logout(ctx context.Context, sessionToken string) error {
	return s.authRepository.DeleteSession(ctx, sessionToken)
}

func (s *service) ChangePassword(ctx context.Context, userID int64, oldPass, newPass string) error {
	current, err := s.authRepository.GetUserPassword(ctx, userID)
	if err != nil {
		return err
	}
	if err := CheckPasswordHash(oldPass, current); err != nil {
		return errors.New("old password is incorrect")
	}
	newHash, err := HashPassword(newPass)
	if err != nil {
		return err
	}
	return s.authRepository.UpdatePassword(ctx, userID, newHash)
}

func (s *service) Register(ctx context.Context, req CreateUserRequest) error {
	email := strings.TrimSpace(req.Email)
	if email == "" || len(req.Password) < 8 {
		return errors.New("invalid email or password too short")
	}
	exists, err := s.authRepository.UserExists(ctx, email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("user already exists")
	}
	hash, err := HashPassword(req.Password)
	if err != nil {
		return err
	}
	return s.authRepository.InsertUser(ctx, email, hash, req.IsAdmin, true)
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return fmt.Sprintf("%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

func CheckPasswordHash(password, encoded string) error {
	parts := strings.Split(encoded, "$")
	if len(parts) != 2 {
		return errors.New("invalid hash format")
	}
	salt, _ := base64.RawStdEncoding.DecodeString(parts[0])
	expectedHash, _ := base64.RawStdEncoding.DecodeString(parts[1])
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	for i := range hash {
		if hash[i] != expectedHash[i] {
			return errors.New("password does not match")
		}
	}
	return nil
}

func GenerateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b), err
}

func (s *service) GetUsers(ctx context.Context) ([]User, error) {
	users, err := s.authRepository.GetUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}
	return users, nil
}
