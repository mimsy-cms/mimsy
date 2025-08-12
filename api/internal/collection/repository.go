package collection

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
	"github.com/mimsy-cms/mimsy/internal/config"
)

type Repository interface {
	FindBySlug(ctx context.Context, slug string) (*Collection, error)
	CollectionExists(ctx context.Context, slug string) (bool, error)
	FindResource(ctx context.Context, collection *Collection, slug string) (*Resource, error)
	FindResources(ctx context.Context, collection *Collection) ([]Resource, error)
	FindAll(ctx context.Context, params *FindAllParams) ([]Collection, error)
	FindAllGlobals(ctx context.Context, params *FindAllParams) ([]Collection, error)
	UpdateResourceContent(ctx context.Context, collection *Collection, resourceSlug string, content map[string]any) (*Resource, error)
	DeleteResource(ctx context.Context, resource *Resource) error
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
	// Id is the private identifier for the resource.
	Id int64
	// Slug is the public identifier for the resource within the collection.
	Slug string
	// CreatedAt is the timestamp when the resource was created.
	CreatedAt time.Time
	// UpdatedAt is the timestamp when the resource was last updated.
	UpdatedAt time.Time
	// Fields is a map of field names to their values.
	Fields map[string]any
	// Collection is the slug of the collection this resource belongs to.
	Collection string
}

// MarshalJSON implements the json.Marshaler interface for Resource.
// We need this custom implementation to handle the conversion of byte slices in the Resource map to JSON.
func (r Resource) MarshalJSON() ([]byte, error) {
	transformed := make(map[string]any)

	transformed["id"] = r.Id
	transformed["slug"] = r.Slug
	transformed["created_at"] = r.CreatedAt
	transformed["updated_at"] = r.UpdatedAt

	for key, value := range r.Fields {
		switch v := value.(type) {
		case []byte:
			var jsonObject any
			if err := json.Unmarshal(v, &jsonObject); err != nil {
				transformed[key] = string(v)
			} else {
				transformed[key] = jsonObject
			}
		default:
			transformed[key] = value
		}
	}

	return json.Marshal(transformed)
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
	fields := map[string]Field{}
	if err := json.Unmarshal(collection.Fields, &fields); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fields: %w", err)
	}

	resource, err := NewSelectQuery(collection.Slug, fields).FindOne(ctx, slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to find resource: %w", err)
	}

	resource.Collection = collection.Slug

	return resource, nil
}

func (r *repository) FindResources(ctx context.Context, collection *Collection) ([]Resource, error) {
	fields := map[string]Field{}
	if err := json.Unmarshal(collection.Fields, &fields); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fields: %w", err)
	}

	resources, err := NewSelectQuery(collection.Slug, fields).FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find resources: %w", err)
	}

	for i := range resources {
		resources[i].Collection = collection.Slug
	}

	return resources, nil
}

type FindAllParams struct {
	Search string
}

func (r *repository) FindAll(ctx context.Context, params *FindAllParams) ([]Collection, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"slug", "name", "fields", "created_at", "created_by", "updated_at", "updated_by", "is_global",
		).
		From("collection").
		Where(sq.Eq{"is_global": false}).
		Where(sq.ILike{`"name"`: fmt.Sprintf("%%%s%%", params.Search)}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	rows, err := config.GetDB(ctx).QueryContext(ctx, sql, args...)
	if err != nil {
		slog.Error("Failed to query collections", "error", err)
		return nil, fmt.Errorf("failed to query collections: %w", err)
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

func (r *repository) FindAllGlobals(ctx context.Context, params *FindAllParams) ([]Collection, error) {
	sql, args, err := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).
		Select(
			"slug", "name", "fields", "created_at", "created_by", "updated_at", "updated_by", "is_global",
		).
		From("collection").
		Where(sq.Eq{"is_global": true}).
		Where(sq.ILike{`"name"`: fmt.Sprintf("%%%s%%", params.Search)}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	rows, err := config.GetDB(ctx).QueryContext(ctx, sql, args...)
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

func (r *repository) DeleteResource(ctx context.Context, resource *Resource) error {
	query := fmt.Sprintf(`DELETE FROM "%s" WHERE slug = $1`, resource.Collection)

	if _, err := config.GetDB(ctx).ExecContext(ctx, query, resource.Slug); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("failed to delete resource: %w", err)
	}

	return nil
}

func (r *repository) UpdateResourceContent(ctx context.Context, collection *Collection, resourceSlug string, content map[string]any) (*Resource, error) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	b := psql.
		Update(pq.QuoteIdentifier(collection.Slug)).
		Where(sq.Eq{"slug": resourceSlug})

	for field, value := range content {
		b = b.Set(pq.QuoteIdentifier(field), value)
	}

	query, args, err := b.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build update SQL query: %w", err)
	}

	if _, err := config.GetDB(ctx).ExecContext(ctx, query, args...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("failed to update resource content: %w", err)
	}

	return r.FindResource(ctx, collection, resourceSlug)
}
