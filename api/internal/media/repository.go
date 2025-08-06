package media

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/mimsy-cms/mimsy/internal/postgres"
)

type Media struct {
	Id           int64
	Uuid         uuid.UUID
	Name         string
	ContentType  string
	CreatedAt    time.Time
	Size         int64
	UploadedById int64
}

type Repository interface {
	Create(ctx context.Context, params *CreateMediaParams) (*Media, error)
	GetById(ctx context.Context, id int64) (*Media, error)
	GetByUuid(ctx context.Context, uuid *uuid.UUID) (*Media, error)
	FindAll(ctx context.Context) ([]Media, error)
	Delete(ctx context.Context, media *Media) error
}

type mediaRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &mediaRepository{db: db}
}

var (
	// ErrMediaReferenced is returned when trying to delete media that is being referenced by other content.
	ErrMediaReferenced = errors.New("media is being referenced by other content")
)

type CreateMediaParams struct {
	Uuid         uuid.UUID
	Name         string
	ContentType  string
	Size         int64
	UploadedById int64
}

func (r *mediaRepository) Create(ctx context.Context, params *CreateMediaParams) (*Media, error) {
	query := `
		INSERT INTO media (uuid, name, content_type, created_at, size, uploaded_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	var id int64

	if err := r.db.QueryRowContext(ctx, query,
		params.Uuid,
		params.Name,
		params.ContentType,
		time.Now(),
		params.Size,
		params.UploadedById,
	).Scan(&id); err != nil {
		return nil, err
	}

	return r.GetById(ctx, id)
}

func (r *mediaRepository) GetById(ctx context.Context, id int64) (*Media, error) {
	query := `SELECT id, uuid, name, content_type, created_at, size, uploaded_by FROM media WHERE id = $1`
	media := &Media{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&media.Id,
		&media.Uuid,
		&media.Name,
		&media.ContentType,
		&media.CreatedAt,
		&media.Size,
		&media.UploadedById)

	if err != nil {
		return nil, err
	}

	return media, nil
}

func (r *mediaRepository) GetByUuid(ctx context.Context, uuid *uuid.UUID) (*Media, error) {
	query := `SELECT id, uuid, name, content_type, created_at, size, uploaded_by FROM media WHERE uuid = $1`
	media := &Media{}

	err := r.db.QueryRowContext(ctx, query, uuid).Scan(
		&media.Id,
		&media.Uuid,
		&media.Name,
		&media.ContentType,
		&media.CreatedAt,
		&media.Size,
		&media.UploadedById)

	if err != nil {
		return nil, err
	}

	return media, nil
}

func (r *mediaRepository) FindAll(ctx context.Context) ([]Media, error) {
	query := `SELECT id, uuid, name, content_type, created_at, size, uploaded_by FROM media`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var medias []Media
	for rows.Next() {
		media := Media{}
		if err := rows.Scan(
			&media.Id,
			&media.Uuid,
			&media.Name,
			&media.ContentType,
			&media.CreatedAt,
			&media.Size,
			&media.UploadedById); err != nil {
			return nil, err
		}
		medias = append(medias, media)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return medias, nil
}

func (r *mediaRepository) Delete(ctx context.Context, media *Media) error {
	query := `DELETE FROM media WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, media.Id)
	if err != nil {
		if postgres.IsErrCode(err, postgres.ErrForeignKeyViolation) {
			return ErrMediaReferenced
		}
		return err
	}

	return nil
}
