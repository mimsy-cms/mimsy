package sync_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mimsy-cms/mimsy/internal/sync"
	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
)

func TestNewSyncStatusRepository(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)
	if repo == nil {
		t.Error("Expected repository to be created")
	}
}

func TestSyncStatusRepository_GetStatus_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"repo", "commit", "commit_message", "commit_date",
		"applied_migration", "applied_at", "is_active", "is_skipped",
		"error_message", "manifest",
	}).AddRow(
		"test-repo", "abc123", "Test commit", now,
		nil, nil, true, false, nil, nil,
	)

	mock.ExpectQuery(`SELECT repo, commit, commit_message, commit_date, applied_migration,
		       applied_at, is_active, is_skipped, error_message, manifest
		FROM sync_status
		WHERE repo = \$1 AND is_active = true
		LIMIT 1`).
		WithArgs("test-repo").
		WillReturnRows(rows)

	status, err := repo.GetStatus("test-repo")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status == nil {
		t.Error("Expected status to not be nil")
	}

	if status.Repo != "test-repo" {
		t.Errorf("Expected repo to be 'test-repo', got %s", status.Repo)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_GetStatus_NoRows(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	mock.ExpectQuery(`SELECT repo, commit, commit_message, commit_date, applied_migration,
		       applied_at, is_active, is_skipped, error_message, manifest
		FROM sync_status
		WHERE repo = \$1 AND is_active = true
		LIMIT 1`).
		WithArgs("test-repo").
		WillReturnError(sql.ErrNoRows)

	status, err := repo.GetStatus("test-repo")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status != nil {
		t.Error("Expected status to be nil")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_GetLastSyncedCommit_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"repo", "commit", "commit_message", "commit_date",
		"applied_migration", "applied_at", "is_active", "is_skipped",
		"error_message", "manifest",
	}).AddRow(
		"test-repo", "abc123", "Test commit", now,
		"{}", now, true, false, nil, "{}",
	)

	mock.ExpectQuery(`SELECT repo, commit, commit_message, commit_date, applied_migration,
		       applied_at, is_active, is_skipped, error_message, manifest
		FROM sync_status
		WHERE repo = \$1 AND applied_at IS NOT NULL
		ORDER BY applied_at DESC
		LIMIT 1`).
		WithArgs("test-repo").
		WillReturnRows(rows)

	status, err := repo.GetLastSyncedCommit("test-repo")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status == nil {
		t.Error("Expected status to not be nil")
	}

	if status.Commit != "abc123" {
		t.Errorf("Expected commit to be 'abc123', got %s", status.Commit)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_MarkError_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	mock.ExpectExec(`UPDATE sync_status
		SET error_message = \$1, is_active = false
		WHERE repo = \$2 AND commit = \$3`).
		WithArgs("test error", "test-repo", "abc123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.MarkError("test-repo", "abc123", errors.New("test error"))
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_CreateIfNotExists_NewRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	commitDate := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM sync_status WHERE repo = \$1 AND commit = \$2`).
		WithArgs("test-repo", "abc123").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectExec(`INSERT INTO sync_status \(repo, commit, commit_message, commit_date, is_active, is_skipped\)
		VALUES \(\$1, \$2, \$3, \$4, false, false\)`).
		WithArgs("test-repo", "abc123", "Test commit", commitDate).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	err = repo.CreateIfNotExists("test-repo", "abc123", "Test commit", commitDate)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_CreateIfNotExists_ExistingRecord(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	commitDate := time.Now()

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT COUNT\(\*\) FROM sync_status WHERE repo = \$1 AND commit = \$2`).
		WithArgs("test-repo", "abc123").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectCommit()

	err = repo.CreateIfNotExists("test-repo", "abc123", "Test commit", commitDate)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_SetManifest_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	schema := mimsy_schema.Schema{
		Collections: []mimsy_schema.Collection{
			{
				Name: "test",
				Schema: mimsy_schema.CollectionFields{
					"title": mimsy_schema.SchemaElement{
						Type: "text",
					},
				},
			},
		},
	}

	expectedJSON, _ := json.Marshal(schema)

	mock.ExpectExec(`UPDATE sync_status
		SET manifest = \$1, applied_at = NOW\(\)
		WHERE repo = \$2 AND commit = \$3`).
		WithArgs(expectedJSON, "test-repo", "abc123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.SetManifest("test-repo", "abc123", schema)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_GetRecentStatuses_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"repo", "commit", "commit_message", "commit_date",
		"applied_migration", "applied_at", "is_active", "is_skipped",
		"error_message", "manifest",
	}).AddRow(
		"test-repo", "abc123", "Test commit", now,
		nil, nil, true, false, nil, nil,
	).AddRow(
		"test-repo", "def456", "Another commit", now.Add(-time.Hour),
		nil, nil, false, false, nil, nil,
	)

	mock.ExpectQuery(`SELECT repo, commit, commit_message, commit_date, applied_migration,
		       applied_at, is_active, is_skipped, error_message, manifest
		FROM sync_status
		ORDER BY commit_date DESC
		LIMIT \$1`).
		WithArgs(5).
		WillReturnRows(rows)

	statuses, err := repo.GetRecentStatuses(5)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(statuses) != 2 {
		t.Errorf("Expected 2 statuses, got %d", len(statuses))
	}

	if statuses[0].Commit != "abc123" {
		t.Errorf("Expected first commit to be 'abc123', got %s", statuses[0].Commit)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_SetAppliedMigration_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	migration := []byte(`{"tables": []}`)

	mock.ExpectExec(`UPDATE sync_status
		SET applied_migration = \$1
		WHERE repo = \$2 AND commit = \$3`).
		WithArgs(migration, "test-repo", "abc123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.SetAppliedMigration("test-repo", "abc123", migration)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_GetActiveMigration_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	now := time.Now()
	rows := sqlmock.NewRows([]string{
		"repo", "commit", "commit_message", "commit_date",
		"applied_migration", "applied_at", "is_active", "is_skipped",
		"error_message", "manifest",
	}).AddRow(
		"test-repo", "abc123", "Test commit", now,
		"{}", now, true, false, nil, "{}",
	)

	mock.ExpectQuery(`SELECT repo, commit, commit_message, commit_date, applied_migration,
									applied_at, is_active, is_skipped, error_message, manifest
		FROM sync_status
		WHERE repo = \$1 AND is_active = true
		LIMIT 1`).
		WithArgs("test-repo").
		WillReturnRows(rows)

	status, err := repo.GetActiveMigration("test-repo")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if status == nil {
		t.Error("Expected status to not be nil")
	}

	if !status.IsActive {
		t.Error("Expected status to be active")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_MarkAsActive_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	mock.ExpectExec(`UPDATE sync_status
		SET is_active = CASE
			WHEN commit = \$2 THEN true
			ELSE false
		END
		WHERE repo = \$1`).
		WithArgs("test-repo", "abc123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.MarkAsActive("test-repo", "abc123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}

func TestSyncStatusRepository_MarkAsSkipped_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock db: %v", err)
	}
	defer db.Close()

	repo := sync.NewSyncStatusRepository(db)

	mock.ExpectExec(`UPDATE sync_status
		SET is_skipped = true, applied_at = NOW\(\)
		WHERE repo = \$1 AND commit = \$2`).
		WithArgs("test-repo", "abc123").
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.MarkAsSkipped("test-repo", "abc123")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %v", err)
	}
}
