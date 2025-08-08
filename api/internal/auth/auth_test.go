package auth_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/mimsy-cms/mimsy/internal/auth"
	mocks "github.com/mimsy-cms/mimsy/internal/mocks/auth"
)

// =================================================================================================
// Helper Functions
// =================================================================================================
func newJSONRequest(t *testing.T, method, url, jsonBody string) *http.Request {
	t.Helper()
	req := httptest.NewRequest(method, url, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func addUserToContext(req *http.Request, user *auth.User) *http.Request {
	ctx := context.WithValue(req.Context(), auth.UserContextKey, user)
	return req.WithContext(ctx)
}

func executeRequest(handler http.Handler, req *http.Request, t *testing.T) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

// =================================================================================================
// HashPassword and CheckPasswordHash
// =================================================================================================

// TestHashPasswordAndCheck tests the HashPassword and CheckPasswordHash functions
func TestHashPasswordAndCheck(t *testing.T) {
	password := "superSecurePassword123!"

	hash, err := auth.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	if hash == "" {
		t.Fatal("hashed password should not be empty")
	}

	if err := auth.CheckPasswordHash(password, hash); err != nil {
		t.Fatalf("password check failed: %v", err)
	}

	wrongPassword := "wrongPassword"
	if err := auth.CheckPasswordHash(wrongPassword, hash); err == nil {
		t.Fatal("expected password check to fail with wrong password, but it succeeded")
	}
}

// =================================================================================================
// generateSessionToken
// =================================================================================================

// TestGenerateSessionToken tests the generateSessionToken function
func TestGenerateSessionToken(t *testing.T) {
	token1, err1 := auth.GenerateSessionToken()
	if err1 != nil {
		t.Fatalf("failed to generate session token1: %v", err1)
	}
	token2, err2 := auth.GenerateSessionToken()
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

// TestLogin_Success tests the login handler for a successful login
func TestLogin_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)

	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	hashedPassword, _ := auth.HashPassword("admin123")

	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), "admin@example.com").Return(&auth.User{
		ID:                 1,
		Email:              "admin@example.com",
		PasswordHash:       hashedPassword,
		MustChangePassword: false,
	}, nil)

	mockRepo.EXPECT().DeleteExpiredSessions(gomock.Any()).Return(nil)

	mockRepo.EXPECT().CreateSession(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	req := newJSONRequest(t, "POST", "/login", `{"email":"admin@example.com","password":"admin123"}`)
	w := executeRequest(http.HandlerFunc(handler.Login), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}
}

// TestLogin_Failure_WrongPassword tests the login handler for a failed login because of incorrect password
func TestLogin_Failure_WrongPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)

	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	hashedPassword, _ := auth.HashPassword("admin123")

	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), "admin@example.com").Return(&auth.User{
		ID:                 1,
		Email:              "admin@example.com",
		PasswordHash:       hashedPassword,
		MustChangePassword: false,
	}, nil)

	req := newJSONRequest(t, "POST", "/login", `{"email":"admin@example.com","password":"wrongpassword"}`)
	w := executeRequest(http.HandlerFunc(handler.Login), req, t)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}
}

// TestLogin_Failure_UserNotFound tests the login handler for a failed login because of user not found
func TestLogin_Failure_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)

	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), "admin@wrongdomain.com").Return(nil, errors.New("user not found"))

	req := newJSONRequest(t, "POST", "/login", `{"email":"admin@wrongdomain.com","password":"admin123"}`)
	w := executeRequest(http.HandlerFunc(handler.Login), req, t)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}
}

// TestLogin_Failure_InvalidRequest tests the login handler for a failed login due to invalid request body
func TestLogin_Failure_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)

	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	// Invalid JSON (missing closing brace)
	req := newJSONRequest(t, "POST", "/login", `{"email":"admin@example.com","password":"admin123"`)

	w := executeRequest(http.HandlerFunc(handler.Login), req, t)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status Bad Request, got %v", w.Code)
	}
}

// TestLogin_Failure_DatabaseError tests the login handler for a failed login due to database error
func TestLogin_Failure_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	mockRepo.EXPECT().
		GetUserByEmail(gomock.Any(), "admin@example.com").
		Return(nil, errors.New("database error"))

	req := newJSONRequest(t, "POST", "/login", `{"email":"admin@example.com","password":"admin123"}`)

	w := executeRequest(http.HandlerFunc(handler.Login), req, t)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}
}

// TestLogin_Failure_SessionCleanupError tests the login handler for a failed login due to session cleanup error
func TestLogin_Failure_SessionCleanupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	hashedPassword, _ := auth.HashPassword("admin123")

	mockRepo.EXPECT().
		GetUserByEmail(gomock.Any(), "admin@example.com").
		Return(&auth.User{
			ID:                 1,
			Email:              "admin@example.com",
			PasswordHash:       hashedPassword,
			MustChangePassword: false,
		}, nil)

	mockRepo.EXPECT().
		DeleteExpiredSessions(gomock.Any()).
		Return(errors.New("session cleanup error"))

	req := newJSONRequest(t, "POST", "/login", `{"email":"admin@example.com","password":"admin123"}`)

	w := executeRequest(http.HandlerFunc(handler.Login), req, t)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// TestLogin_Failure_SessionInsertError tests the login handler for a failed login due to session insert error
func TestLogin_Failure_SessionInsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	hashedPassword, _ := auth.HashPassword("admin123")

	mockRepo.EXPECT().
		GetUserByEmail(gomock.Any(), "admin@example.com").
		Return(&auth.User{
			ID:                 1,
			Email:              "admin@example.com",
			PasswordHash:       hashedPassword,
			MustChangePassword: false,
		}, nil)

	mockRepo.EXPECT().
		DeleteExpiredSessions(gomock.Any()).
		Return(nil)

	mockRepo.EXPECT().
		CreateSession(gomock.Any(), gomock.Any(), int64(1), gomock.Any()).
		Return(errors.New("session insert error"))

	req := newJSONRequest(t, "POST", "/login", `{"email":"admin@example.com","password":"admin123"}`)

	w := executeRequest(http.HandlerFunc(handler.Login), req, t)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// =================================================================================================
// LogoutHandler
// =================================================================================================

// TestLogout_Success tests the logout handler
func TestLogout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	sessionID := "test-session-id"

	mockRepo.EXPECT().
		DeleteSession(gomock.Any(), sessionID).
		Return(nil)

	req := newJSONRequest(t, "POST", "/logout", "")
	req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})

	w := executeRequest(http.HandlerFunc(handler.Logout), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "Logged out successfully") {
		t.Errorf("expected response body to contain 'Logged out successfully', got: %s", w.Body.String())
	}
}

// TestLogout_Failure_NoSessionCookie tests the logout handler when no session cookie is present
func TestLogout_Failure_NoSessionCookie(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	req := httptest.NewRequest("POST", "/logout", nil)
	w := executeRequest(http.HandlerFunc(handler.Logout), req, t)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "No session found") {
		t.Errorf("expected response body to contain 'No session found', got: %s", w.Body.String())
	}
}

// TestLogout_Failure_DatabaseError tests the logout handler for a failed logout due to database error
func TestLogout_Failure_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	sessionID := "test-session-id"

	mockRepo.EXPECT().
		DeleteSession(gomock.Any(), sessionID).
		Return(errors.New("database error"))

	req := newJSONRequest(t, "POST", "/logout", "")
	req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})

	w := executeRequest(http.HandlerFunc(handler.Logout), req, t)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// =================================================================================================
// ChangePasswordHandler
// =================================================================================================

// TestChangePassword_Success tests the change password handler for a successful password change
func TestChangePassword_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	user := &auth.User{
		ID:                 1,
		Email:              "admin@example.com",
		IsAdmin:            true,
		MustChangePassword: false,
	}

	hashedPassword, _ := auth.HashPassword("admin123")

	mockRepo.EXPECT().
		GetUserPassword(gomock.Any(), user.ID).
		Return(hashedPassword, nil)

	mockRepo.EXPECT().
		UpdatePassword(gomock.Any(), user.ID, gomock.Any()).
		Return(nil)

	reqBody := `{"old_password":"admin123","new_password":"newpassword"}`

	req := newJSONRequest(t, "POST", "/password", reqBody)
	req = addUserToContext(req, user)

	w := executeRequest(http.HandlerFunc(handler.ChangePassword), req, t)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status OK, got %v", w.Code)
	}
}

// TestChangePassword_Failure_WrongOldPassword tests the change password handler for a failed password change because of incorrect old password
func TestChangePassword_Failure_WrongOldPassword(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	user := &auth.User{
		ID:                 1,
		Email:              "admin@example.com",
		IsAdmin:            true,
		MustChangePassword: false,
	}

	hashedPassword, err := auth.HashPassword("admin123")
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	mockRepo.EXPECT().
		GetUserPassword(gomock.Any(), user.ID).
		Return(hashedPassword, nil)

	reqBody := `{"old_password":"wrongpassword","new_password":"newpassword"}`

	req := newJSONRequest(t, "POST", "/password", reqBody)
	req = addUserToContext(req, user)

	w := executeRequest(http.HandlerFunc(handler.ChangePassword), req, t)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "old password is incorrect") {
		t.Errorf("expected response body to contain 'old password is incorrect', got: %s", w.Body.String())
	}
}

// TestChangePassword_Failure_InvalidRequest tests the change password handler for a failed password change due to invalid request body
func TestChangePassword_Failure_InvalidRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	user := &auth.User{
		ID:                 1,
		Email:              "admin@example.com",
		IsAdmin:            true,
		MustChangePassword: false,
	}

	req := newJSONRequest(t, "POST", "/password", `{"old_password":"admin123","new_password"`) // malformed JSON
	req = addUserToContext(req, user)

	w := executeRequest(http.HandlerFunc(handler.ChangePassword), req, t)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status Bad Request, got %v", w.Code)
	}
}

// TestChangePassword_Failure_DatabaseError tests the change password handler for a failed password change due to database error
func TestChangePassword_Failure_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	user := &auth.User{
		ID:                 1,
		Email:              "admin@example.com",
		IsAdmin:            true,
		MustChangePassword: false,
	}

	hashedPassword, _ := auth.HashPassword("admin123")

	mockRepo.EXPECT().
		GetUserPassword(gomock.Any(), user.ID).
		Return(hashedPassword, nil)

	mockRepo.EXPECT().
		UpdatePassword(gomock.Any(), user.ID, gomock.Any()).
		Return(errors.New("database error"))

	req := newJSONRequest(t, "POST", "/password", `{"old_password":"admin123","new_password":"newpassword"}`)
	req = addUserToContext(req, user)

	w := executeRequest(http.HandlerFunc(handler.ChangePassword), req, t)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
	}
}

// TestChangePassword_Failure_MissingUser tests the change password handler for a failed password change due to missing user in context
func TestChangePassword_Failure_MissingUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	req := newJSONRequest(t, "POST", "/password", `{"old_password":"admin123","new_password":"newpassword"}`)

	w := executeRequest(http.HandlerFunc(handler.ChangePassword), req, t)

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

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)

	mockRepo.EXPECT().
		CountUsers(gomock.Any()).
		Return(0, nil)

	mockRepo.EXPECT().
		InsertUser(gomock.Any(), gomock.Any(), gomock.Any(), true, true).
		Return(nil)

	err := authService.CreateAdminUser(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

// TestCreateAdminUser_Failure_UserCountError tests the CreateAdminUser function for a failed admin user creation due to user count error
func TestCreateAdminUser_Failure_UserCountError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)

	mockRepo.EXPECT().
		CountUsers(gomock.Any()).
		Return(0, errors.New("database error"))

	err := authService.CreateAdminUser(context.Background())
	if err == nil || !strings.Contains(err.Error(), "failed to count users") {
		t.Fatalf("expected error about counting users, got %v", err)
	}
}

// TestCreateAdminUser_Failure_UserInsertError tests the CreateAdminUser function for a failed admin user creation due to user insert error
func TestCreateAdminUser_Failure_UserInsertError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)

	mockRepo.EXPECT().
		CountUsers(gomock.Any()).
		Return(0, nil)

	mockRepo.EXPECT().
		InsertUser(gomock.Any(), gomock.Any(), gomock.Any(), true, true).
		Return(errors.New("insert error"))

	err := authService.CreateAdminUser(context.Background())
	if err == nil || !strings.Contains(err.Error(), "failed to insert admin user") {
		t.Fatalf("expected error about creating admin user, got %v", err)
	}
}

// TestCreateAdminUser_Failure_UserAlreadyExists tests the CreateAdminUser function for a failed admin user creation due to user already exists
func TestCreateAdminUser_Failure_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)

	mockRepo.EXPECT().
		CountUsers(gomock.Any()).
		Return(1, nil)

	err := authService.CreateAdminUser(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
