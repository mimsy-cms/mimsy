package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mockauth "github.com/mimsy-cms/mimsy/internal/mocks/auth"
)

// TestHashPasswordAndCheck tests the HashPassword and CheckPasswordHash functions
func TestHashPasswordAndCheck(t *testing.T) {
	password := "superSecurePassword123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash == "" {
		t.Fatal("hashed password should not be empty")
	}

	if err := CheckPasswordHash(password, hash); err != nil {
		t.Fatalf("password check failed: %v", err)
	}

	wrongPassword := "wrongPassword"
	if err := CheckPasswordHash(wrongPassword, hash); err == nil {
		t.Fatal("expected password check to fail with wrong password, but it succeeded")
	}
}

// TestGenerateSalt tests the generateSalt function
func TestGenerateSalt(t *testing.T) {
	saltLen := 16
	salt1, err1 := generateSalt(saltLen)
	if err1 != nil {
		t.Fatalf("failed to generate salt1: %v", err1)
	}
	salt2, err2 := generateSalt(saltLen)
	if err2 != nil {
		t.Fatalf("failed to generate salt2: %v", err2)
	}

	if len(salt1) == 0 {
		t.Fatal("generated salt should not be empty")
	}
	if string(salt1) == string(salt2) {
		t.Fatal("generated salts should not be the same")
	}
}

// TestCompareHashes tests the compareHashes function
func TestCompareHashes(t *testing.T) {
	tests := []struct {
		name     string
		hash1    string
		hash2    string
		expected bool
	}{
		{"same hashes", "abc123", "abc123", true},
		{"different hashes", "abc123", "xyz789", false},
		{"different lengths", "abc123", "abc1234", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareHashes([]byte(tt.hash1), []byte(tt.hash2))
			if result != tt.expected {
				t.Errorf("compareHashes(%q, %q) = %v; want %v", tt.hash1, tt.hash2, result, tt.expected)
			}
		})
	}
}

// TestGenerateSessionToken tests the generateSessionToken function
func TestGenerateSessionToken(t *testing.T) {
	token1, err1 := generateSessionToken()
	if err1 != nil {
		t.Fatalf("failed to generate session token1: %v", err1)
	}
	token2, err2 := generateSessionToken()
	if err2 != nil {
		t.Fatalf("failed to generate session token2: %v", err2)
	}

	if len(token1) == 0 || len(token2) == 0 {
		t.Fatal("generated session token should not be empty")
	}
	if token1 == token2 {
		t.Fatal("generated session tokens should not be the same")
	}
}

// TestLoginHandler_Success tests the login handler for a successful login
func TestLoginHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)
	mockRow := mockauth.NewMockRow(ctrl)

	mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(dest ...interface{}) error {
			*dest[0].(*int64) = int64(1)
			*dest[1].(*string) = "admin@example.com"
			hash, _ := HashPassword("admin123")
			*dest[2].(*string) = hash
			*dest[3].(*bool) = false
			return nil
		},
	)

	mockDB.EXPECT().QueryRow(`SELECT id, email, password, must_change_password FROM "user" WHERE email = $1`, "admin@example.com").Return(mockRow)

	mockDB.EXPECT().Exec(gomock.Any(), gomock.Any()).Return(nil, nil)

	mockDB.EXPECT().Exec(
		gomock.Any(),
		gomock.Any(),
		int64(1),
		gomock.Any(),
	).Return(nil, nil)

	handler := LoginHandler(mockDB)

	body := strings.NewReader(`{"email":"admin@example.com","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}
}

// TestLoginHandler_Failure tests the login handler for a failed login
func TestLoginHandler_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)
	mockRow := mockauth.NewMockRow(ctrl)

	mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
		func(dest ...interface{}) error {
			*dest[0].(*int64) = int64(1)
			*dest[1].(*string) = "admin@example.com"
			hash, _ := HashPassword("wrongpassword")
			*dest[2].(*string) = hash
			*dest[3].(*bool) = false
			return nil
		},
	)

	mockDB.EXPECT().QueryRow(`SELECT id, email, password, must_change_password FROM "user" WHERE email = $1`, "admin@example.com").Return(mockRow)

	handler := LoginHandler(mockDB)

	body := strings.NewReader(`{"email":"admin@example.com","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}
}
