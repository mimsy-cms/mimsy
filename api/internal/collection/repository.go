package collection

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/mimsy-cms/mimsy/internal/config"
)

type Repository interface {
	FindBySlug(ctx context.Context, slug string) (*Collection, error)
	CollectionExists(ctx context.Context, slug string) (bool, error)
	FindResource(ctx context.Context, collection *Collection, slug string) (*Resource, error)
	FindResources(ctx context.Context, collection *Collection) ([]Resource, error)
	FindAll(ctx context.Context) ([]Collection, error)
	ListGlobals(ctx context.Context) ([]Collection, error)
	UpdateResourceContent(ctx context.Context, collection *Collection, resourceSlug string, content map[string]any) (*Resource, error)
}

type repository struct{}

func NewRepository() *repository {
	return &repository{}
}

type Collection struct {
	Slug      string
	Name      string
	Fields    json.RawMessage
	CreatedAt string
	CreatedBy string
	UpdatedAt string
	UpdatedBy *string
	IsGlobal  bool
}

type Resource struct {
	ID        int             `json:"id"`
	Slug      string          `json:"slug"`
	Content   json.RawMessage `json:"content"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// MarshalJSON implements the json.Marshaler interface for Resource.
// We need this custom implementation to handle the conversion of byte slices in the Resource map to JSON.
func (r Resource) MarshalJSON() ([]byte, error) {
	var contentData map[string]any
	if len(r.Content) > 0 {
		_ = json.Unmarshal(r.Content, &contentData)
	} else {
		contentData = make(map[string]any)
	}

	response := map[string]any{
		"id":         r.ID,
		"slug":       r.Slug,
		"created_at": r.CreatedAt,
		"updated_at": r.UpdatedAt,
		"content":    contentData,
	}

	for key, value := range contentData {
		response[key] = value
	}

	return json.Marshal(response)
}

type Field struct {
	Type       string         `json:"type"`
	Label      string         `json:"label"`
	Required   bool           `json:"required,omitempty"`
	Default    any            `json:"default,omitempty"`
	Options    []string       `json:"options,omitempty"`
	Relation   *FieldRelation `json:"relation,omitempty"`
	Validation *Validation    `json:"validation,omitempty"`
}

type Validation struct {
	MinLength int    `json:"min_length,omitempty"`
	MaxLength int    `json:"max_length,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
	Min       *int   `json:"min,omitempty"`
	Max       *int   `json:"max,omitempty"`
}

type FieldRelationType string

const (
	FieldRelationTypeManyToOne  FieldRelationType = "many-to-one"
	FieldRelationTypeManyToMany FieldRelationType = "many-to-many"
)

type FieldRelation struct {
	To   string
	Type FieldRelationType
}

var ErrNotFound = errors.New("not found")

func (r *repository) FindBySlug(ctx context.Context, slug string) (*Collection, error) {
	var collection Collection
	err := config.GetDB(ctx).QueryRowContext(ctx,
		`SELECT slug, name, fields, created_at, created_by, updated_at, updated_by FROM "collection" WHERE slug = $1`,
		slug,
	).Scan(
		&collection.Slug,
		&collection.Name,
		&collection.Fields,
		&collection.CreatedAt,
		&collection.CreatedBy,
		&collection.UpdatedAt,
		&collection.UpdatedBy,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &collection, nil
}

func (r *repository) CollectionExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := config.GetDB(ctx).QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM "collection" WHERE slug = $1)`, slug).Scan(&exists)
	return exists, err
}

func (r *repository) FindResource(ctx context.Context, collection *Collection, slug string) (*Resource, error) {
	query := fmt.Sprintf(`SELECT id, slug, content, created_at, updated_at FROM %s WHERE slug = $1`, pq.QuoteIdentifier(collection.Slug))

	var resource Resource
	err := config.GetDB(ctx).QueryRowContext(ctx, query, slug).Scan(
		&resource.ID,
		&resource.Slug,
		&resource.Content,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, err
	}

	return &resource, nil
}

func (r *repository) FindResources(ctx context.Context, collection *Collection) ([]Resource, error) {
	query := fmt.Sprintf(`SELECT id, slug, content, created_at, updated_at FROM %s`, pq.QuoteIdentifier(collection.Slug))

	rows, err := config.GetDB(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query resources: %w", err)
	}
	defer rows.Close()

	var resources []Resource
	for rows.Next() {
		var resource Resource
		if err := rows.Scan(&resource.ID, &resource.Slug, &resource.Content, &resource.CreatedAt, &resource.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan resource: %w", err)
		}
		resources = append(resources, resource)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over resources: %w", err)
	}

	return resources, nil
}

func (r *repository) FindAll(ctx context.Context) ([]Collection, error) {
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

func (r *repository) ListGlobals(ctx context.Context) ([]Collection, error) {
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

func (r *repository) UpdateResourceContent(ctx context.Context, collection *Collection, resourceSlug string, content map[string]any) (*Resource, error) {
	contentJSON, err := json.Marshal(content)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal content: %w", err)
	}

	query := fmt.Sprintf(
		`UPDATE %s SET content = $1, updated_at = NOW() WHERE slug = $2 RETURNING id, slug, content, created_at, updated_at`,
		pq.QuoteIdentifier(collection.Slug),
	)

	var resource Resource
	err = config.GetDB(ctx).QueryRowContext(ctx, query, contentJSON, resourceSlug).Scan(
		&resource.ID,
		&resource.Slug,
		&resource.Content,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, fmt.Errorf("failed to update resource content: %w", err)
	}

	return &resource, nil
}
