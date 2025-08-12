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
	ErrorMessage     string    `json:"error_message"`
}

type SyncStatusRepository interface {
	GetStatus(repo string) (*SyncStatus, error)
	GetLastSyncedCommit(repo string) (*SyncStatus, error)
	GetRecentStatuses(limit int) ([]*SyncStatus, error)
	MarkError(repo string, commitSha string, err error) error
	CreateStatus(repo string, commitSha string, commitMessage string, commitDate time.Time) error
	SetManifest(repo string, commitSha string, manifest mimsy_schema.Schema) error
	SetAppliedMigration(repo string, commitSha string, migration []byte) error
}

type syncStatusRepository struct {
	db *sql.DB
}

func NewSyncStatusRepository(db *sql.DB) SyncStatusRepository {
	return &syncStatusRepository{db: db}
}

func (r *syncStatusRepository) GetStatus(repo string) (*SyncStatus, error) {
	query := `
		SELECT repo, commit, commit_message, commit_date, applied_migration,
		       applied_at, is_active, error_message, manifest
		FROM sync_status
		WHERE repo = $1 AND is_active = true
		LIMIT 1`

	var status SyncStatus
	var appliedMigration, manifest sql.NullString
	var appliedAt sql.NullTime

	err := r.db.QueryRow(query, repo).Scan(
		&status.Repo,
		&status.Commit,
		&status.CommitMessage,
		&status.CommitDate,
		&appliedMigration,
		&appliedAt,
		&status.IsActive,
		&status.ErrorMessage,
		&manifest,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get sync status: %w", err)
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

	return &status, nil
}

func (r *syncStatusRepository) GetLastSyncedCommit(repo string) (*SyncStatus, error) {
	query := `
		SELECT repo, commit, commit_message, commit_date, applied_migration,
		       applied_at, is_active, error_message, manifest
		FROM sync_status
		WHERE repo = $1 AND applied_at IS NOT NULL
		ORDER BY applied_at DESC
		LIMIT 1`

	var status SyncStatus
	var appliedMigration, manifest sql.NullString
	var appliedAt sql.NullTime

	err := r.db.QueryRow(query, repo).Scan(
		&status.Repo,
		&status.Commit,
		&status.CommitMessage,
		&status.CommitDate,
		&appliedMigration,
		&appliedAt,
		&status.IsActive,
		&status.ErrorMessage,
		&manifest,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get last synced commit: %w", err)
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

	return &status, nil
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

func (r *syncStatusRepository) CreateStatus(repo string, commitSha string, commitMessage string, commitDate time.Time) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Deactivate previous active status
	_, err = tx.Exec("UPDATE sync_status SET is_active = false WHERE repo = $1", repo)
	if err != nil {
		return fmt.Errorf("failed to deactivate previous status: %w", err)
	}

	// Create new active status
	query := `
		INSERT INTO sync_status (repo, commit, commit_message, commit_date, is_active)
		VALUES ($1, $2, $3, $4, true)`

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

func (r *syncStatusRepository) GetRecentStatuses(limit int) ([]*SyncStatus, error) {
	query := `
		SELECT repo, commit, commit_message, commit_date, applied_migration,
		       applied_at, is_active, error_message, manifest
		FROM sync_status
		ORDER BY commit_date DESC
		LIMIT $1`

	rows, err := r.db.Query(query, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent statuses: %w", err)
	}
	defer rows.Close()

	var statuses []*SyncStatus
	for rows.Next() {
		var status SyncStatus
		var appliedMigration, manifest sql.NullString
		var appliedAt sql.NullTime

		err := rows.Scan(
			&status.Repo,
			&status.Commit,
			&status.CommitMessage,
			&status.CommitDate,
			&appliedMigration,
			&appliedAt,
			&status.IsActive,
			&status.ErrorMessage,
			&manifest,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan status row: %w", err)
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

		statuses = append(statuses, &status)
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
