package collection

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/mimsy-cms/mimsy/internal/config"
)

type Repository interface {
	FindBySlug(ctx context.Context, slug string) (*Collection, error)
	CollectionExists(ctx context.Context, slug string) (bool, error)
	FindItemsBySlug(ctx context.Context, slug string) ([]Item, error)
	List(ctx context.Context) ([]Collection, error)
	ListGlobals(ctx context.Context) ([]Collection, error)
}

type PostgresRepository struct{}

func NewRepository() *PostgresRepository {
	return &PostgresRepository{}
}

type Collection struct {
	Slug      string
	Name      string
	Fields    []byte
	CreatedAt string
	CreatedBy string
	UpdatedAt string
	UpdatedBy *string
	IsGlobal  bool
}

type Item struct {
	ID           int             `json:"id"`
	ResourceSlug string          `json:"slug"`
	Data         json.RawMessage `json:"data"`
}

var ErrNotFound = errors.New("not found")

func (r *PostgresRepository) FindBySlug(ctx context.Context, slug string) (*Collection, error) {
	var coll Collection
	err := config.GetDB(ctx).QueryRowContext(ctx,
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
	err := config.GetDB(ctx).QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM "collection" WHERE slug = $1)`, slug).Scan(&exists)
	return exists, err
}

func (r *PostgresRepository) FindItemsBySlug(ctx context.Context, slug string) ([]Item, error) {
	rows, err := config.GetDB(ctx).QueryContext(ctx, `SELECT id, data, slug FROM "collection_item" WHERE collection_slug = $1`, slug)
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

func (r *PostgresRepository) List(ctx context.Context) ([]Collection, error) {
	rows, err := config.GetDB(ctx).QueryContext(ctx, `SELECT slug, name, fields, created_at, created_by, updated_at, updated_by, is_global FROM "collection" WHERE is_global = false`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var coll Collection
		if err := rows.Scan(&coll.Slug, &coll.Name, &coll.Fields, &coll.CreatedAt, &coll.CreatedBy, &coll.UpdatedAt, &coll.UpdatedBy, &coll.IsGlobal); err != nil {
			return nil, err
		}
		collections = append(collections, coll)
	}
	return collections, nil
}

func (r *PostgresRepository) ListGlobals(ctx context.Context) ([]Collection, error) {
	rows, err := config.GetDB(ctx).QueryContext(ctx, `SELECT slug, name, fields, created_at, created_by, updated_at, updated_by, is_global FROM "collection" WHERE is_global = true`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var coll Collection
		if err := rows.Scan(&coll.Slug, &coll.Name, &coll.Fields, &coll.CreatedAt, &coll.CreatedBy, &coll.UpdatedAt, &coll.UpdatedBy, &coll.IsGlobal); err != nil {
			return nil, err
		}
		collections = append(collections, coll)
	}
	return collections, nil
}
