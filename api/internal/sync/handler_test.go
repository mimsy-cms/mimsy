package sync_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/mimsy-cms/mimsy/internal/auth"
	"github.com/mimsy-cms/mimsy/internal/cron"
	mocks_cron "github.com/mimsy-cms/mimsy/internal/mocks/cron"
	mocks_sync "github.com/mimsy-cms/mimsy/internal/mocks/sync"
	"github.com/mimsy-cms/mimsy/internal/sync"
)

// Helper function to create authenticated request
func addUserToContext(req *http.Request) *http.Request {
	user := &auth.User{
		ID:    1,
		Email: "test@example.com",
	}
	ctx := context.WithValue(req.Context(), auth.UserContextKey, user)
	return req.WithContext(ctx)
}

func TestHandler_Status_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	handler := sync.NewHandler(mockRepo, mockCron)

	now := time.Now()
	expectedStatuses := []sync.SyncStatus{
		{
			Repo:          "test-repo",
			Commit:        "abc123",
			CommitMessage: "Test commit",
			CommitDate:    now,
			IsActive:      true,
			AppliedAt:     now,
		},
	}

	mockRepo.EXPECT().
		GetRecentStatuses(5).
		Return(expectedStatuses, nil).
		Times(1)

	req := httptest.NewRequest("GET", "/sync/status", nil)
	req = addUserToContext(req)
	w := httptest.NewRecorder()

	handler.Status(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}
}

func TestHandler_Status_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	handler := sync.NewHandler(mockRepo, mockCron)

	req := httptest.NewRequest("GET", "/sync/status", nil)
	w := httptest.NewRecorder()

	handler.Status(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandler_Status_CustomLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	handler := sync.NewHandler(mockRepo, mockCron)

	mockRepo.EXPECT().
		GetRecentStatuses(3).
		Return([]sync.SyncStatus{}, nil).
		Times(1)

	req := httptest.NewRequest("GET", "/sync/status?limit=3", nil)
	req = addUserToContext(req)
	w := httptest.NewRecorder()

	handler.Status(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandler_Status_InvalidLimit(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	handler := sync.NewHandler(mockRepo, mockCron)

	// Should use default limit of 5 when limit is too high
	mockRepo.EXPECT().
		GetRecentStatuses(5).
		Return([]sync.SyncStatus{}, nil).
		Times(1)

	req := httptest.NewRequest("GET", "/sync/status?limit=15", nil)
	req = addUserToContext(req)
	w := httptest.NewRecorder()

	handler.Status(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandler_Status_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	handler := sync.NewHandler(mockRepo, mockCron)

	mockRepo.EXPECT().
		GetRecentStatuses(5).
		Return(nil, errors.New("database error")).
		Times(1)

	req := httptest.NewRequest("GET", "/sync/status", nil)
	req = addUserToContext(req)
	w := httptest.NewRecorder()

	handler.Status(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestHandler_Jobs_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	handler := sync.NewHandler(mockRepo, mockCron)

	expectedJobs := []cron.JobStatus{
		{
			Name:      "test-job",
			Schedule:  "*/5 * * * *",
			IsRunning: true,
		},
	}

	mockCron.EXPECT().
		GetJobStatuses(gomock.Any()).
		Return(expectedJobs, nil).
		Times(1)

	req := httptest.NewRequest("GET", "/sync/jobs", nil)
	req = addUserToContext(req)
	w := httptest.NewRecorder()

	handler.Jobs(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHandler_Jobs_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	handler := sync.NewHandler(mockRepo, mockCron)

	req := httptest.NewRequest("GET", "/sync/jobs", nil)
	w := httptest.NewRecorder()

	handler.Jobs(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestHandler_Jobs_CronServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	handler := sync.NewHandler(mockRepo, mockCron)

	mockCron.EXPECT().
		GetJobStatuses(gomock.Any()).
		Return(nil, errors.New("cron service error")).
		Times(1)

	req := httptest.NewRequest("GET", "/sync/jobs", nil)
	req = addUserToContext(req)
	w := httptest.NewRecorder()

	handler.Jobs(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestNewHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks_sync.NewMockSyncStatusRepository(ctrl)
	mockCron := mocks_cron.NewMockCronService(ctrl)

	handler := sync.NewHandler(mockRepo, mockCron)

	if handler == nil {
		t.Error("expected handler to be created")
	}
}
