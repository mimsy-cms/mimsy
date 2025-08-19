package collection

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/lib/pq"

	sq "github.com/Masterminds/squirrel"
	"github.com/mimsy-cms/mimsy/internal/config"
	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
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
	FindUserEmail(ctx context.Context, id int64) (string, error)
}

type repository struct{}

func NewRepository() *repository {
	return &repository{}
}

type Collection struct {
	Slug           string
	Name           string
	Fields         json.RawMessage
	CreatedAt      string
	CreatedBy      int64
	CreatedByEmail string
	UpdatedAt      string
	UpdatedBy      int64
	UpdatedByEmail string
	IsGlobal       bool
}

type Resource struct {
	// Id is the private identifier for the resource.
	Id int64
	// Slug is the public identifier for the resource within the collection.
	Slug string
	// CreatedAt is the timestamp when the resource was created.
	CreatedAt time.Time
	// CreatedBy is the identifier of the user who created the resource.
	CreatedBy int64
	// CreatedByEmail is the email address of the user who created the resource.
	CreatedByEmail string
	// UpdatedAt is the timestamp when the resource was last updated.
	UpdatedAt time.Time
	// UpdatedBy is the identifier of the user who last updated the resource.
	UpdatedBy int64
	// UpdatedByEmail is the email address of the user who last updated the resource.
	UpdatedByEmail string
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
	transformed["created_by_email"] = r.CreatedByEmail
	transformed["updated_at"] = r.UpdatedAt
	transformed["updated_by"] = r.UpdatedBy
	transformed["updated_by_email"] = r.UpdatedByEmail

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

	collection.CreatedByEmail, _ = r.FindUserEmail(ctx, collection.CreatedBy)
	collection.UpdatedByEmail, _ = r.FindUserEmail(ctx, collection.UpdatedBy)

	return &collection, nil
}

func (r *repository) CollectionExists(ctx context.Context, slug string) (bool, error) {
	var exists bool
	err := config.GetDB(ctx).QueryRowContext(ctx, `SELECT EXISTS(SELECT 1 FROM "collection" WHERE slug = $1)`, slug).Scan(&exists)
	return exists, err
}

func (r *repository) FindResource(ctx context.Context, collection *Collection, slug string) (*Resource, error) {
	fields := mimsy_schema.CollectionFields{}
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
	fields := mimsy_schema.CollectionFields{}
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
		Where(sq.Eq{"slug": resourceSlug}).
		Set("updated_at", sq.Expr("NOW()"))

	if updatedBy, exists := content["updated_by"]; exists {
		b = b.Set("updated_by", updatedBy)
	}

	for field, value := range content {
		// Skip read only columns that should not be updated
		if slices.Contains(readOnlyColumns, field) {
			continue
		}

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

func (r *repository) FindUserEmail(ctx context.Context, id int64) (string, error) {
	var email string
	err := config.GetDB(ctx).QueryRowContext(ctx,
		`SELECT email FROM "user" WHERE id = $1`,
		id,
	).Scan(&email)

	if err == sql.ErrNoRows {
		return "", ErrNotFound
	} else if err != nil {
		return "", fmt.Errorf("failed to find user email: %w", err)
	}

	return email, nil
}
