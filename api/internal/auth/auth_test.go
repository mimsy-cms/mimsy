package auth

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	mockauth "github.com/mimsy-cms/mimsy/internal/mocks/auth"
)

// =================================================================================================
// HashPassword and CheckPasswordHash
// =================================================================================================

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

// =================================================================================================
// generateSalt
// =================================================================================================

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

// =================================================================================================
// compareHashes
// =================================================================================================

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

// =================================================================================================
// generateSessionToken
// =================================================================================================

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

// =================================================================================================
// LoginHandler
// =================================================================================================

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

// TestLoginHandler_Failure_WrongPassword tests the login handler for a failed login because of incorrect password
func TestLoginHandler_Failure_WrongPassword(t *testing.T) {
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

	handler := LoginHandler(mockDB)

	body := strings.NewReader(`{"email":"admin@example.com","password":"wrongpassword"}`)
	req := httptest.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}
}

// TestLoginHandler_Failure_UserNotFound tests the login handler for a failed login because of user not found
func TestLoginHandler_Failure_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)
	mockRow := mockauth.NewMockRow(ctrl)

	mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(sql.ErrNoRows)

	mockDB.EXPECT().QueryRow(`SELECT id, email, password, must_change_password FROM "user" WHERE email = $1`, "admin@wrongdomain.com").Return(mockRow)

	handler := LoginHandler(mockDB)

	body := strings.NewReader(`{"email":"admin@wrongdomain.com","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}
}

// TestLoginHandler_Failure_InvalidRequest tests the login handler for a failed login due to invalid request body
func TestLoginHandler_Failure_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)

	handler := LoginHandler(mockDB)

	body := strings.NewReader(`{"email":"admin@example.com","password":"admin123"`)
	req := httptest.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status Bad Request, got %v", w.Code)
	}
}

// TestLoginHandler_Failure_DatabaseError tests the login handler for a failed login due to database error
func TestLoginHandler_Failure_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)
	mockRow := mockauth.NewMockRow(ctrl)

	mockDB.EXPECT().QueryRow(`SELECT id, email, password, must_change_password FROM "user" WHERE email = $1`, "admin@example.com").Return(mockRow)

	mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("database error"))

	handler := LoginHandler(mockDB)

	body := strings.NewReader(`{"email":"admin@example.com","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// TestLoginHandler_Failure_SessionCleanupError tests the login handler for a failed login due to session cleanup error
func TestLoginHandler_Failure_SessionCleanupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)
	mockRow := mockauth.NewMockRow(ctrl)

	mockDB.EXPECT().QueryRow(`SELECT id, email, password, must_change_password FROM "user" WHERE email = $1`, "admin@example.com").Return(mockRow)

	mockRow.EXPECT().
		Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(args ...any) error {
			*(args[0].(*int64)) = 1
			*(args[1].(*string)) = "admin@example.com"
			hash, _ := HashPassword("admin123")
			*(args[2].(*string)) = hash
			*(args[3].(*bool)) = false
			return nil
		})

	mockDB.EXPECT().Exec(`DELETE FROM "session" WHERE expires_at < NOW()`).Return(nil, errors.New("session cleanup error"))

	handler := LoginHandler(mockDB)

	body := strings.NewReader(`{"email":"admin@example.com","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// TestLoginHandler_Failure_SessionInsertError tests the login handler for a failed login due to session insert error
func TestLoginHandler_Failure_SessionInsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)
	mockRow := mockauth.NewMockRow(ctrl)

	mockDB.EXPECT().QueryRow(`SELECT id, email, password, must_change_password FROM "user" WHERE email = $1`, "admin@example.com").Return(mockRow)

	mockRow.EXPECT().
		Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		DoAndReturn(func(args ...any) error {
			*(args[0].(*int64)) = 1
			*(args[1].(*string)) = "admin@example.com"
			hash, _ := HashPassword("admin123")
			*(args[2].(*string)) = hash
			*(args[3].(*bool)) = false
			return nil
		})

	mockDB.EXPECT().Exec(`DELETE FROM "session" WHERE expires_at < NOW()`).Return(nil, nil)

	mockDB.EXPECT().Exec(
		`INSERT INTO session (id, user_id, expires_at) VALUES ($1, $2, $3)`,
		gomock.Any(),
		int64(1),
		gomock.Any(),
	).Return(nil, errors.New("session insert error"))

	handler := LoginHandler(mockDB)

	body := strings.NewReader(`{"email":"admin@example.com","password":"admin123"}`)
	req := httptest.NewRequest("POST", "/login", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// =================================================================================================
// LogoutHandler
// =================================================================================================

// TestLogoutHandler_Success tests the logout handler
func TestLogoutHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)

	sessionID := "test-session-id"

	mockDB.EXPECT().Exec(`DELETE FROM session WHERE id = $1`, sessionID).Return(nil, nil)

	handler := LogoutHandler(mockDB)

	req := httptest.NewRequest("POST", "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Logged out successfully") {
		t.Errorf("expected response body to contain 'Logged out successfully', got: %s", w.Body.String())
	}
}

// TestLogoutHandler_Failure tests the logout handler when no session cookie is present
func TestLogoutHandler_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)

	handler := LogoutHandler(mockDB)

	req := httptest.NewRequest("POST", "/logout", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "No session found") {
		t.Errorf("expected response body to contain 'No session found', got: %s", w.Body.String())
	}
}

// TestLogoutHandler_Failure_DatabaseError tests the logout handler for a failed logout due to database error
func TestLogoutHandler_Failure_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)

	sessionID := "test-session-id"

	mockDB.EXPECT().Exec(`DELETE FROM session WHERE id = $1`, sessionID).Return(nil, errors.New("database error"))

	handler := LogoutHandler(mockDB)

	req := httptest.NewRequest("POST", "/logout", nil)
	req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// =================================================================================================
// ChangePasswordHandler
// =================================================================================================

// TestChangePasswordHandler_Success tests the change password handler for a successful password change
func TestChangePasswordHandler_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)
	mockRow := mockauth.NewMockRow(ctrl)

	user := &User{
		ID:                 1,
		Email:              "admin@example.com",
		IsAdmin:            true,
		MustChangePassword: false,
	}

	hashedPassword, _ := HashPassword("admin123")
	user.PasswordHash = hashedPassword

	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(
		func(dest ...interface{}) error {
			*dest[0].(*string) = user.PasswordHash
			return nil
		},
	)

	mockDB.EXPECT().QueryRow(`SELECT password FROM "user" WHERE id = $1`, user.ID).Return(mockRow)

	mockDB.EXPECT().Exec(
		`UPDATE "user" SET password = $1, must_change_password = FALSE WHERE id = $2`,
		gomock.Any(), user.ID,
	).Return(nil, nil)

	handler := ChangePasswordHandler(mockDB)

	body := strings.NewReader(`{"old_password":"admin123","new_password":"newpassword"}`)
	req := httptest.NewRequest("POST", "/password", body)
	req.AddCookie(&http.Cookie{Name: "session", Value: "1"})
	req.Header.Set("Content-Type", "application/json")

	req = addUserToContext(req, user)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}
}

func addUserToContext(req *http.Request, user *User) *http.Request {
	ctx := context.WithValue(req.Context(), userContextKey, user)
	return req.WithContext(ctx)
}

// TestChangePasswordHandler_Failure_WrongOldPassword tests the change password handler for a failed password change because of incorrect old password
func TestChangePasswordHandler_Failure_WrongOldPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)
	mockRow := mockauth.NewMockRow(ctrl)

	user := &User{
		ID:                 1,
		Email:              "admin@example.com",
		IsAdmin:            true,
		MustChangePassword: false,
	}

	hashedPassword, err := HashPassword("admin123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}
	user.PasswordHash = hashedPassword

	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(
		func(dest ...interface{}) error {
			*dest[0].(*string) = user.PasswordHash
			return nil
		},
	)

	mockDB.EXPECT().QueryRow(`SELECT password FROM "user" WHERE id = $1`, user.ID).Return(mockRow)

	handler := ChangePasswordHandler(mockDB)

	body := strings.NewReader(`{"old_password":"wrongpassword","new_password":"newpassword"}`)
	req := httptest.NewRequest("POST", "/password", body)
	req.Header.Set("Content-Type", "application/json")

	req = addUserToContext(req, user)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Old password is incorrect") {
		t.Errorf("expected response body to contain 'Old password is incorrect', got: %s", w.Body.String())
	}
}

// TestChangePasswordHandler_Failure_InvalidRequest tests the change password handler for a failed password change due to invalid request body
func TestChangePasswordHandler_Failure_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)

	handler := ChangePasswordHandler(mockDB)

	body := strings.NewReader(`{"old_password":"admin123","new_password"`)
	req := httptest.NewRequest("POST", "/password", body)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "session", Value: "1"})

	user := &User{
		ID:                 1,
		Email:              "admin@example.com",
		IsAdmin:            true,
		MustChangePassword: false,
	}
	req = addUserToContext(req, user)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status Bad Request, got %v", w.Code)
	}
}

// TestChangePasswordHandler_Failure_DatabaseError tests the change password handler for a failed password change due to database error
func TestChangePasswordHandler_Failure_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)
	mockRow := mockauth.NewMockRow(ctrl)

	user := &User{
		ID:                 1,
		Email:              "admin@example.com",
		IsAdmin:            true,
		MustChangePassword: false,
	}

	hashedPassword, _ := HashPassword("admin123")
	user.PasswordHash = hashedPassword

	mockDB.EXPECT().QueryRow(`SELECT password FROM "user" WHERE id = $1`, user.ID).Return(mockRow)

	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(
		func(dest ...interface{}) error {
			*dest[0].(*string) = user.PasswordHash
			return nil
		},
	)

	mockDB.EXPECT().Exec(
		`UPDATE "user" SET password = $1, must_change_password = FALSE WHERE id = $2`,
		gomock.Any(), user.ID,
	).Return(nil, errors.New("database error"))

	handler := ChangePasswordHandler(mockDB)

	body := strings.NewReader(`{"old_password":"admin123","new_password":"newpassword"}`)
	req := httptest.NewRequest("POST", "/password", body)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: "session", Value: "1"})

	req = addUserToContext(req, user)

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// TestChangePasswordHandler_Failure_MissingUser tests the change password handler for a failed password change due to missing user in context
func TestChangePasswordHandler_Failure_MissingUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)

	handler := ChangePasswordHandler(mockDB)

	body := strings.NewReader(`{"old_password":"admin123","new_password":"newpassword"}`)
	req := httptest.NewRequest("POST", "/password", body)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}
}

// =================================================================================================
// CreateAdminUser
// =================================================================================================

// TestCreateAdminUser_Success tests the CreateAdminUser function for a successful admin user creation
func TestCreateAdminUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mockauth.NewMockDB(ctrl)

	mockDB.EXPECT().QueryRowContext(gomock.Any(), `SELECT COUNT(*) FROM "user"`).Return(mockauth.NewMockRow(ctrl)).DoAndReturn(
		func(ctx context.Context, query string, args ...interface{}) *mockauth.MockRow {
			row := mockauth.NewMockRow(ctrl)
			row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
				*dest[0].(*int) = 0
				return nil
			})
			return row
		},
	)

	mockDB.EXPECT().ExecContext(gomock.Any(), `INSERT INTO "user" (email, password, must_change_password, is_admin) VALUES ($1, $2, $3, $4)`,
		gomock.Any(), gomock.Any(), true, true,
	).Return(nil, nil)

	err := CreateAdminUser(context.Background(), mockDB)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestCreateAdminUser_Failure_UserCountError tests the CreateAdminUser function for a failed admin user creation due to user count error

// TestCreateAdminUser_Failure_UserInsertError tests the CreateAdminUser function for a failed admin user creation due to user insert error

// TestCreateAdminUser_Failure_UserAlreadyExists tests the CreateAdminUser function for a failed admin user creation due to user already exists
