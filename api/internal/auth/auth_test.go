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
	auth_interface "github.com/mimsy-cms/mimsy/internal/interfaces/auth"
	mocks "github.com/mimsy-cms/mimsy/internal/mocks/auth"
)

// =================================================================================================
// Helper Functions
// =================================================================================================
func setupMocks(t *testing.T) (*gomock.Controller, *mocks.MockDB, *mocks.MockRow) {
	ctrl := gomock.NewController(t)
	mockDB := mocks.NewMockDB(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	t.Cleanup(func() {
		ctrl.Finish()
	})
	return ctrl, mockDB, mockRow
}

func newJSONRequest(t *testing.T, method, url, jsonBody string) *http.Request {
	req := httptest.NewRequest(method, url, strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func addUserToContext(req *http.Request, user *auth_interface.User) *http.Request {
	ctx := context.WithValue(req.Context(), auth.UserContextKey, user)
	return req.WithContext(ctx)
}

func expectUserQuery(mockDB *mocks.MockDB, mockRow *mocks.MockRow, email, password string) {
	mockDB.EXPECT().QueryRow(`SELECT id, email, password, must_change_password FROM "user" WHERE email = $1`, email).Return(mockRow)
	mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
		*dest[0].(*int64) = int64(1)
		*dest[1].(*string) = email
		*dest[2].(*string) = password
		*dest[3].(*bool) = false
		return nil
	})
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

	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), "admin@example.com").Return(&auth_interface.User{
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

	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), "admin@example.com").Return(&auth_interface.User{
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

// TestLoginHandler_Failure_UserNotFound tests the login handler for a failed login because of user not found
func TestLoginHandler_Failure_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)

	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	// Simulate user not found
	mockRepo.EXPECT().GetUserByEmail(gomock.Any(), "admin@wrongdomain.com").Return(nil, errors.New("user not found"))

	// DeleteExpiredSessions and CreateSession should NOT be called

	req := newJSONRequest(t, "POST", "/login", `{"email":"admin@wrongdomain.com","password":"admin123"}`)
	w := executeRequest(http.HandlerFunc(handler.Login), req, t)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}
}

// TestLoginHandler_Failure_InvalidRequest tests the login handler for a failed login due to invalid request body
func TestLoginHandler_Failure_InvalidRequest(t *testing.T) {
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

// TestLoginHandler_Failure_DatabaseError tests the login handler for a failed login due to database error
func TestLoginHandler_Failure_DatabaseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	// Simulate DB error during GetUserByEmail
	mockRepo.EXPECT().
		GetUserByEmail(gomock.Any(), "admin@example.com").
		Return(nil, errors.New("database error"))

	req := newJSONRequest(t, "POST", "/login", `{"email":"admin@example.com","password":"admin123"}`)

	w := executeRequest(http.HandlerFunc(handler.Login), req, t)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status Unauthorized, got %v", w.Code)
	}
}

// TestLoginHandler_Failure_SessionCleanupError tests the login handler for a failed login due to session cleanup error
func TestLoginHandler_Failure_SessionCleanupError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockAuthRepository(ctrl)
	authService := auth.NewAuthService(mockRepo)
	handler := auth.NewHandler(authService)

	hashedPassword, _ := auth.HashPassword("admin123")

	mockRepo.EXPECT().
		GetUserByEmail(gomock.Any(), "admin@example.com").
		Return(&auth_interface.User{
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

// // TestLoginHandler_Failure_SessionInsertError tests the login handler for a failed login due to session insert error
// func TestLoginHandler_Failure_SessionInsertError(t *testing.T) {
// 	_, mockDB, mockRow := setupMocks(t)

// 	hashedPassword, _ := HashPassword("admin123")
// 	expectUserQuery(mockDB, mockRow, "admin@example.com", hashedPassword)

// 	mockDB.EXPECT().Exec(`DELETE FROM "session" WHERE expires_at < NOW()`).Return(nil, nil)

// 	mockDB.EXPECT().Exec(
// 		`INSERT INTO session (id, user_id, expires_at) VALUES ($1, $2, $3)`,
// 		gomock.Any(),
// 		int64(1),
// 		gomock.Any(),
// 	).Return(nil, errors.New("session insert error"))

// 	handler := LoginHandler(mockDB)

// 	req := newJSONRequest(t, "POST", "/login", `{"email":"admin@example.com","password":"admin123"}`)

// 	w := executeRequest(handler, req, t)

// 	if w.Code != http.StatusInternalServerError {
// 		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
// 	}
// }

// // =================================================================================================
// // LogoutHandler
// // =================================================================================================

// // TestLogoutHandler_Success tests the logout handler
// func TestLogoutHandler_Success(t *testing.T) {
// 	_, mockDB, _ := setupMocks(t)

// 	sessionID := "test-session-id"

// 	mockDB.EXPECT().Exec(`DELETE FROM session WHERE id = $1`, sessionID).Return(nil, nil)

// 	handler := LogoutHandler(mockDB)

// 	req := newJSONRequest(t, "POST", "/logout", "")
// 	req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})

// 	w := executeRequest(handler, req, t)

// 	if w.Code != http.StatusOK {
// 		t.Fatalf("expected status OK, got %v", w.Code)
// 	}

// 	if !strings.Contains(w.Body.String(), "Logged out successfully") {
// 		t.Errorf("expected response body to contain 'Logged out successfully', got: %s", w.Body.String())
// 	}
// }

// // TestLogoutHandler_Failure tests the logout handler when no session cookie is present
// func TestLogoutHandler_Failure(t *testing.T) {
// 	_, mockDB, _ := setupMocks(t)

// 	handler := LogoutHandler(mockDB)

// 	req := httptest.NewRequest("POST", "/logout", nil)
// 	w := executeRequest(handler, req, t)

// 	if w.Code != http.StatusUnauthorized {
// 		t.Fatalf("expected status Unauthorized, got %v", w.Code)
// 	}

// 	if !strings.Contains(w.Body.String(), "No session found") {
// 		t.Errorf("expected response body to contain 'No session found', got: %s", w.Body.String())
// 	}
// }

// // TestLogoutHandler_Failure_DatabaseError tests the logout handler for a failed logout due to database error
// func TestLogoutHandler_Failure_DatabaseError(t *testing.T) {
// 	_, mockDB, _ := setupMocks(t)

// 	sessionID := "test-session-id"

// 	mockDB.EXPECT().Exec(`DELETE FROM session WHERE id = $1`, sessionID).Return(nil, errors.New("database error"))

// 	handler := LogoutHandler(mockDB)

// 	req := newJSONRequest(t, "POST", "/logout", "")
// 	req.AddCookie(&http.Cookie{Name: "session", Value: sessionID})

// 	w := executeRequest(handler, req, t)

// 	if w.Code != http.StatusInternalServerError {
// 		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
// 	}
// }

// // =================================================================================================
// // ChangePasswordHandler
// // =================================================================================================

// // TestChangePasswordHandler_Success tests the change password handler for a successful password change
// func TestChangePasswordHandler_Success(t *testing.T) {
// 	_, mockDB, mockRow := setupMocks(t)

// 	user := &User{
// 		ID:                 1,
// 		Email:              "admin@example.com",
// 		IsAdmin:            true,
// 		MustChangePassword: false,
// 	}

// 	hashedPassword, _ := HashPassword("admin123")
// 	user.PasswordHash = hashedPassword

// 	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(
// 		func(dest ...interface{}) error {
// 			*dest[0].(*string) = user.PasswordHash
// 			return nil
// 		},
// 	)

// 	mockDB.EXPECT().QueryRow(`SELECT password FROM "user" WHERE id = $1`, user.ID).Return(mockRow)

// 	mockDB.EXPECT().Exec(
// 		`UPDATE "user" SET password = $1, must_change_password = FALSE WHERE id = $2`,
// 		gomock.Any(), user.ID,
// 	).Return(nil, nil)

// 	handler := ChangePasswordHandler(mockDB)

// 	req := newJSONRequest(t, "POST", "/password", `{"old_password":"admin123","new_password":"newpassword"}`)
// 	req.AddCookie(&http.Cookie{Name: "session", Value: "1"})
// 	req = addUserToContext(req, user)

// 	w := executeRequest(handler, req, t)

// 	if w.Code != http.StatusOK {
// 		t.Fatalf("expected status OK, got %v", w.Code)
// 	}
// }

// // TestChangePasswordHandler_Failure_WrongOldPassword tests the change password handler for a failed password change because of incorrect old password
// func TestChangePasswordHandler_Failure_WrongOldPassword(t *testing.T) {
// 	_, mockDB, mockRow := setupMocks(t)

// 	user := &User{
// 		ID:                 1,
// 		Email:              "admin@example.com",
// 		IsAdmin:            true,
// 		MustChangePassword: false,
// 	}

// 	hashedPassword, err := HashPassword("admin123")
// 	if err != nil {
// 		t.Fatalf("failed to hash password: %v", err)
// 	}
// 	user.PasswordHash = hashedPassword

// 	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(
// 		func(dest ...interface{}) error {
// 			*dest[0].(*string) = user.PasswordHash
// 			return nil
// 		},
// 	)

// 	mockDB.EXPECT().QueryRow(`SELECT password FROM "user" WHERE id = $1`, user.ID).Return(mockRow)

// 	handler := ChangePasswordHandler(mockDB)

// 	req := newJSONRequest(t, "POST", "/password", `{"old_password":"wrongpassword","new_password":"newpassword"}`)

// 	req = addUserToContext(req, user)

// 	w := executeRequest(handler, req, t)

// 	if w.Code != http.StatusUnauthorized {
// 		t.Fatalf("expected status Unauthorized, got %v", w.Code)
// 	}

// 	if !strings.Contains(w.Body.String(), "Old password is incorrect") {
// 		t.Errorf("expected response body to contain 'Old password is incorrect', got: %s", w.Body.String())
// 	}
// }

// // TestChangePasswordHandler_Failure_InvalidRequest tests the change password handler for a failed password change due to invalid request body
// func TestChangePasswordHandler_Failure_InvalidRequest(t *testing.T) {
// 	_, mockDB, _ := setupMocks(t)

// 	handler := ChangePasswordHandler(mockDB)

// 	req := newJSONRequest(t, "POST", "/password", `{"old_password":"admin123","new_password"`)
// 	req.AddCookie(&http.Cookie{Name: "session", Value: "1"})

// 	user := &User{
// 		ID:                 1,
// 		Email:              "admin@example.com",
// 		IsAdmin:            true,
// 		MustChangePassword: false,
// 	}
// 	req = addUserToContext(req, user)

// 	w := httptest.NewRecorder()
// 	handler.ServeHTTP(w, req)

// 	if w.Code != http.StatusBadRequest {
// 		t.Fatalf("expected status Bad Request, got %v", w.Code)
// 	}
// }

// // TestChangePasswordHandler_Failure_DatabaseError tests the change password handler for a failed password change due to database error
// func TestChangePasswordHandler_Failure_DatabaseError(t *testing.T) {
// 	_, mockDB, mockRow := setupMocks(t)

// 	user := &User{
// 		ID:                 1,
// 		Email:              "admin@example.com",
// 		IsAdmin:            true,
// 		MustChangePassword: false,
// 	}

// 	hashedPassword, _ := HashPassword("admin123")
// 	user.PasswordHash = hashedPassword

// 	mockDB.EXPECT().QueryRow(`SELECT password FROM "user" WHERE id = $1`, user.ID).Return(mockRow)

// 	mockRow.EXPECT().Scan(gomock.Any()).DoAndReturn(
// 		func(dest ...interface{}) error {
// 			*dest[0].(*string) = user.PasswordHash
// 			return nil
// 		},
// 	)

// 	mockDB.EXPECT().Exec(
// 		`UPDATE "user" SET password = $1, must_change_password = FALSE WHERE id = $2`,
// 		gomock.Any(), user.ID,
// 	).Return(nil, errors.New("database error"))

// 	handler := ChangePasswordHandler(mockDB)

// 	req := newJSONRequest(t, "POST", "/password", `{"old_password":"admin123","new_password":"newpassword"}`)
// 	req = addUserToContext(req, user)

// 	w := executeRequest(handler, req, t)

// 	if w.Code != http.StatusInternalServerError {
// 		t.Fatalf("expected status Internal Server Error, got %v", w.Code)
// 	}
// }

// // TestChangePasswordHandler_Failure_MissingUser tests the change password handler for a failed password change due to missing user in context
// func TestChangePasswordHandler_Failure_MissingUser(t *testing.T) {
// 	_, mockDB, _ := setupMocks(t)

// 	handler := ChangePasswordHandler(mockDB)

// 	req := newJSONRequest(t, "POST", "/password", `{"old_password":"admin123","new_password":"newpassword"}`)

// 	w := executeRequest(handler, req, t)

// 	if w.Code != http.StatusUnauthorized {
// 		t.Fatalf("expected status Unauthorized, got %v", w.Code)
// 	}
// }

// // =================================================================================================
// // CreateAdminUser
// // =================================================================================================

// // TestCreateAdminUser_Success tests the CreateAdminUser function for a successful admin user creation
// func TestCreateAdminUser_Success(t *testing.T) {
// 	ctrl, mockDB, _ := setupMocks(t)

// 	mockDB.EXPECT().QueryRowContext(gomock.Any(), `SELECT COUNT(*) FROM "user"`).Return(mockauth.NewMockRow(ctrl)).DoAndReturn(
// 		func(ctx context.Context, query string, args ...interface{}) *mockauth.MockRow {
// 			row := mockauth.NewMockRow(ctrl)
// 			row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
// 				*dest[0].(*int) = 0
// 				return nil
// 			})
// 			return row
// 		},
// 	)

// 	mockDB.EXPECT().ExecContext(gomock.Any(), `INSERT INTO "user" (email, password, must_change_password, is_admin) VALUES ($1, $2, $3, $4)`,
// 		gomock.Any(), gomock.Any(), true, true,
// 	).Return(nil, nil)

// 	err := CreateAdminUser(context.Background(), mockDB)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}
// }

// // TestCreateAdminUser_Failure_UserCountError tests the CreateAdminUser function for a failed admin user creation due to user count error
// func TestCreateAdminUser_Failure_UserCountError(t *testing.T) {
// 	ctrl, mockDB, _ := setupMocks(t)

// 	mockDB.EXPECT().QueryRowContext(gomock.Any(), `SELECT COUNT(*) FROM "user"`).DoAndReturn(
// 		func(ctx context.Context, query string, args ...interface{}) auth_interface.Row {
// 			row := mockauth.NewMockRow(ctrl)
// 			row.EXPECT().Scan(gomock.Any()).Return(errors.New("database error"))
// 			return row
// 		},
// 	)

// 	err := CreateAdminUser(context.Background(), mockDB)
// 	if err == nil || !strings.Contains(err.Error(), "failed to count users") {
// 		t.Fatalf("expected error about counting users, got %v", err)
// 	}
// }

// // TestCreateAdminUser_Failure_UserInsertError tests the CreateAdminUser function for a failed admin user creation due to user insert error
// func TestCreateAdminUser_Failure_UserInsertError(t *testing.T) {
// 	ctrl, mockDB, _ := setupMocks(t)

// 	mockDB.EXPECT().QueryRowContext(gomock.Any(), `SELECT COUNT(*) FROM "user"`).DoAndReturn(
// 		func(ctx context.Context, query string, args ...interface{}) auth_interface.Row {
// 			row := mockauth.NewMockRow(ctrl)
// 			row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
// 				*dest[0].(*int) = 0
// 				return nil
// 			})
// 			return row
// 		},
// 	)

// 	mockDB.EXPECT().ExecContext(gomock.Any(), `INSERT INTO "user" (email, password, must_change_password, is_admin) VALUES ($1, $2, $3, $4)`,
// 		gomock.Any(), gomock.Any(), true, true,
// 	).Return(nil, errors.New("insert error"))

// 	err := CreateAdminUser(context.Background(), mockDB)
// 	if err == nil || !strings.Contains(err.Error(), "failed to create admin user") {
// 		t.Fatalf("expected error about creating admin user, got %v", err)
// 	}
// }

// // TestCreateAdminUser_Failure_UserAlreadyExists tests the CreateAdminUser function for a failed admin user creation due to user already exists
// func TestCreateAdminUser_Failure_UserAlreadyExists(t *testing.T) {
// 	ctrl, mockDB, _ := setupMocks(t)

// 	mockDB.EXPECT().QueryRowContext(gomock.Any(), `SELECT COUNT(*) FROM "user"`).DoAndReturn(
// 		func(ctx context.Context, query string, args ...interface{}) auth_interface.Row {
// 			row := mockauth.NewMockRow(ctrl)
// 			row.EXPECT().Scan(gomock.Any()).DoAndReturn(func(dest ...interface{}) error {
// 				*dest[0].(*int) = 1 // Simulate that a user already exists
// 				return nil
// 			})
// 			return row
// 		},
// 	)

// 	err := CreateAdminUser(context.Background(), mockDB)
// 	if err != nil {
// 		t.Fatalf("expected no error, got %v", err)
// 	}
// }
