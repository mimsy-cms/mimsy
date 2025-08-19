package sync

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
)

type SyncStatus struct {
	Repo             string    `json:"repo"`
	Commit           string    `json:"commit"`
	CommitMessage    string    `json:"commit_message"`
	CommitDate       time.Time `json:"commit_date"`
	Manifest         string    `json:"manifest"`
	AppliedMigration string    `json:"applied_migration"`
	AppliedAt        time.Time `json:"applied_at"`
	IsActive         bool      `json:"is_active"`
	IsSkipped        bool      `json:"is_skipped"`
	ErrorMessage     string    `json:"error_message"`
}

type SyncStatusRepository interface {
	GetStatus(repo string) (*SyncStatus, error)
	GetLastSyncedCommit(repo string) (*SyncStatus, error)
	GetRecentStatuses(limit int) ([]SyncStatus, error)
	MarkError(repo string, commitSha string, err error) error
	CreateIfNotExists(repo string, commitSha string, commitMessage string, commitDate time.Time) error
	SetManifest(repo string, commitSha string, manifest mimsy_schema.Schema) error
	SetAppliedMigration(repo string, commitSha string, migration []byte) error
	GetActiveMigration(repo string) (*SyncStatus, error)
	MarkAsActive(repo string, commitSha string) error
	MarkAsSkipped(repo string, commitSha string) error
}

type syncStatusRepository struct {
	db *sql.DB
}

func NewSyncStatusRepository(db *sql.DB) SyncStatusRepository {
	return &syncStatusRepository{db: db}
}

// scanSyncStatus is a helper function to scan database rows into SyncStatus struct
func scanSyncStatus(scanner interface {
	Scan(dest ...interface{}) error
}) (*SyncStatus, error) {
	var status SyncStatus
	var appliedMigration, manifest, errorMessage sql.NullString
	var appliedAt sql.NullTime

	err := scanner.Scan(
		&status.Repo,
		&status.Commit,
		&status.CommitMessage,
		&status.CommitDate,
		&appliedMigration,
		&appliedAt,
		&status.IsActive,
		&status.IsSkipped,
		&errorMessage,
		&manifest,
	)

	if err != nil {
		return nil, err
	}

	if appliedMigration.Valid {
		status.AppliedMigration = appliedMigration.String
	}

	if appliedAt.Valid {
		status.AppliedAt = appliedAt.Time
	}

	if manifest.Valid {
		status.Manifest = manifest.String
	}

	if errorMessage.Valid {
		status.ErrorMessage = errorMessage.String
	}

	return &status, nil
}

func (r *syncStatusRepository) GetStatus(repo string) (*SyncStatus, error) {
	query := `
		SELECT repo, commit, commit_message, commit_date, applied_migration,
		       applied_at, is_active, is_skipped, error_message, manifest
		FROM sync_status
		WHERE repo = $1 AND is_active = true
		LIMIT 1`

	row := r.db.QueryRow(query, repo)
	status, err := scanSyncStatus(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sync status: %w", err)
	}

	return status, nil
}

func (r *syncStatusRepository) GetLastSyncedCommit(repo string) (*SyncStatus, error) {
	query := `
		SELECT repo, commit, commit_message, commit_date, applied_migration,
		       applied_at, is_active, is_skipped, error_message, manifest
		FROM sync_status
		WHERE repo = $1 AND applied_at IS NOT NULL
		ORDER BY applied_at DESC
		LIMIT 1`

	row := r.db.QueryRow(query, repo)
	status, err := scanSyncStatus(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last synced commit: %w", err)
	}

	return status, nil
}

func (r *syncStatusRepository) MarkError(repo string, commitSha string, err error) error {
	query := `
		UPDATE sync_status
		SET error_message = $1, is_active = false
		WHERE repo = $2 AND commit = $3`

	_, execErr := r.db.Exec(query, err.Error(), repo, commitSha)
	if execErr != nil {
		return fmt.Errorf("failed to mark error: %w", execErr)
	}

	return nil
}

func (r *syncStatusRepository) CreateIfNotExists(repo string, commitSha string, commitMessage string, commitDate time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Check if the (repo, commitSha) pair already exists
	var count int
	err = tx.QueryRow("SELECT COUNT(*) FROM sync_status WHERE repo = $1 AND commit = $2", repo, commitSha).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing status: %w", err)
	}

	// If it already exists, do nothing
	if count > 0 {
		return tx.Commit()
	}

	// Create new status
	query := `
		INSERT INTO sync_status (repo, commit, commit_message, commit_date, is_active, is_skipped)
		VALUES ($1, $2, $3, $4, false, false)`

	_, err = tx.Exec(query, repo, commitSha, commitMessage, commitDate)
	if err != nil {
		return fmt.Errorf("failed to create sync status: %w", err)
	}

	return tx.Commit()
}

func (r *syncStatusRepository) SetManifest(repo string, commitSha string, manifest mimsy_schema.Schema) error {
	manifestJSON, err := json.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	query := `
		UPDATE sync_status
		SET manifest = $1, applied_at = NOW()
		WHERE repo = $2 AND commit = $3`

	_, err = r.db.Exec(query, manifestJSON, repo, commitSha)
	if err != nil {
		return fmt.Errorf("failed to set manifest: %w", err)
	}

	return nil
}

func (r *syncStatusRepository) GetRecentStatuses(limit int) ([]SyncStatus, error) {
	query := `
		SELECT repo, commit, commit_message, commit_date, applied_migration,
		       applied_at, is_active, is_skipped, error_message, manifest
		FROM sync_status
		ORDER BY commit_date DESC
		LIMIT $1`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent statuses: %w", err)
	}
	defer rows.Close()

	var statuses []SyncStatus
	for rows.Next() {
		status, err := scanSyncStatus(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan status row: %w", err)
		}

		statuses = append(statuses, *status)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating status rows: %w", err)
	}

	return statuses, nil
}

func (r *syncStatusRepository) SetAppliedMigration(repo string, commitSha string, migration []byte) error {
	query := `
		UPDATE sync_status
		SET applied_migration = $1
		WHERE repo = $2 AND commit = $3`

	_, err := r.db.Exec(query, migration, repo, commitSha)
	if err != nil {
		return fmt.Errorf("failed to set applied migration: %w", err)
	}

	return nil
}

func (r *syncStatusRepository) GetActiveMigration(repo string) (*SyncStatus, error) {
	query := `
		SELECT repo, commit, commit_message, commit_date, applied_migration,
									applied_at, is_active, is_skipped, error_message, manifest
		FROM sync_status
		WHERE repo = $1 AND is_active = true
		LIMIT 1`

	row := r.db.QueryRow(query, repo)
	status, err := scanSyncStatus(row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get active migration: %w", err)
	}

	return status, nil
}

func (s *syncStatusRepository) MarkAsActive(repo string, commitSha string) error {
	query := `
		UPDATE sync_status
		SET is_active = CASE
			WHEN commit = $2 THEN true
			ELSE false
		END
		WHERE repo = $1`

	_, err := s.db.Exec(query, repo, commitSha)
	if err != nil {
		return fmt.Errorf("failed to mark as active: %w", err)
	}

	return nil
}

func (s *syncStatusRepository) MarkAsSkipped(repo string, commitSha string) error {
	query := `
		UPDATE sync_status
		SET is_skipped = true, applied_at = NOW()
		WHERE repo = $1 AND commit = $2`

	_, err := s.db.Exec(query, repo, commitSha)
	if err != nil {
		return fmt.Errorf("failed to mark as skipped: %w", err)
	}

	return nil
}
