package collection

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

type Repository interface {
	FindBySlug(ctx context.Context, slug string) (*Collection, error)
	CollectionExists(ctx context.Context, slug string) (bool, error)
	FindItemsBySlug(ctx context.Context, slug string) ([]Item, error)
}

type PostgresRepository struct {
	DB *sql.DB
}

func NewRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{DB: db}
}

type Collection struct {
	Name      string
	Fields    []byte
	CreatedAt string
	CreatedBy string
	UpdatedAt string
	UpdatedBy string
}

type Item struct {
	ID           int             `json:"id"`
	ResourceSlug string          `json:"slug"`
	Data         json.RawMessage `json:"data"`
}

var ErrNotFound = errors.New("not found")

func (r *PostgresRepository) FindBySlug(ctx context.Context, slug string) (*Collection, error) {
	var coll Collection
	err := r.DB.QueryRowContext(ctx,
		`SELECT name, fields, created_at, created_by, updated_at, updated_by FROM "collection" WHERE slug = $1`,
		slug,
	).Scan(&coll.Name, &coll.Fields, &coll.CreatedAt, &coll.CreatedBy, &coll.UpdatedAt, &coll.UpdatedBy)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return &coll, nil
}

func (r *PostgresRepository) CollectionExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := r.DB.QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM "collection" WHERE slug = $1)`, slug).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) FindItemsBySlug(ctx context.Context, slug string) ([]Item, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, data, slug FROM "collection_item" WHERE collection_slug = $1`, slug)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		if err := rows.Scan(&item.ID, &item.Data, &item.ResourceSlug); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}
