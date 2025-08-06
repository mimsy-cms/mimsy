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

type AuthService interface {
	CreateAdminUser(ctx context.Context) error
	Login(ctx context.Context, email, password string) (*LoginResponse, error)
	Logout(ctx context.Context, sessionToken string) error
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
	Register(ctx context.Context, req CreateUserRequest) error
}

type authService struct {
	repo AuthRepository
}

func NewAuthService(repo AuthRepository) AuthService {
	return &authService{repo}
}

func (s *authService) CreateAdminUser(ctx context.Context) error {
	count, err := s.repo.CountUsers(ctx)
	if err != nil || count > 0 {
		return err
	}
	hash, err := HashPassword(adminPassword)
	if err != nil {
		return err
	}
	return s.repo.InsertUser(ctx, adminEmail, hash, true, true)
}

func (s *authService) Login(ctx context.Context, email, password string) (*LoginResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("user not found")
	}
	if err := CheckPasswordHash(password, user.PasswordHash); err != nil {
		return nil, errors.New("invalid credentials")
	}
	if err := s.repo.DeleteExpiredSessions(ctx); err != nil {
		return nil, err
	}
	token, err := generateSessionToken()
	if err != nil {
		return nil, err
	}
	expiresAt := time.Now().Add(7 * 24 * time.Hour)
	if err := s.repo.CreateSession(ctx, token, user.ID, expiresAt); err != nil {
		return nil, err
	}
	return &LoginResponse{Session: token, MustChangePassword: user.MustChangePassword}, nil
}

func (s *authService) Logout(ctx context.Context, sessionToken string) error {
	return s.repo.DeleteSession(ctx, sessionToken)
}

func (s *authService) ChangePassword(ctx context.Context, userID int64, oldPass, newPass string) error {
	current, err := s.repo.GetUserPassword(ctx, userID)
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
	return s.repo.UpdatePassword(ctx, userID, newHash)
}

func (s *authService) Register(ctx context.Context, req CreateUserRequest) error {
	email := strings.TrimSpace(req.Email)
	if email == "" || len(req.Password) < 8 {
		return errors.New("invalid email or password too short")
	}
	exists, err := s.repo.UserExists(ctx, email)
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
	return s.repo.InsertUser(ctx, email, hash, req.IsAdmin, true)
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

func generateSessionToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b), err
}
