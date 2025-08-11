package collection

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/mimsy-cms/mimsy/internal/config"
)

type Repository interface {
	FindBySlug(ctx context.Context, slug string) (*Collection, error)
	CollectionExists(ctx context.Context, slug string) (bool, error)
	FindResources(ctx context.Context, collection *Collection) ([]Resource, error)
	List(ctx context.Context) ([]Collection, error)
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
}

type Resource map[string]any

type Field struct {
	Type     string
	Relation *FieldRelation
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

func (r *repository) FindResources(ctx context.Context, collection *Collection) ([]Resource, error) {
	fields := map[string]Field{}
	if err := json.Unmarshal(collection.Fields, &fields); err != nil {
		return nil, fmt.Errorf("failed to unmarshal fields: %w", err)
	}

	queryFields := []string{"id", "slug"}
	for name, field := range fields {
		if field.Type == "relation" {
			if field.Relation.Type == FieldRelationTypeManyToOne {
				queryFields = append(queryFields, fmt.Sprintf("%s_id", name))
			}
		} else {
			queryFields = append(queryFields, name)
		}
	}

	quotedQueryFields := make([]string, len(queryFields))
	for i, field := range queryFields {
		quotedQueryFields[i] = pq.QuoteIdentifier(field)
	}

	query := fmt.Sprintf(
		`SELECT %s FROM %s`,
		strings.Join(quotedQueryFields, ", "),
		pq.QuoteIdentifier(collection.Slug),
	)

	rows, err := config.GetDB(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// We need to dynamically scan the rows as we dynamically built the selected columns.
	values := make([]any, len(queryFields))
	valuesPtrs := make([]any, len(queryFields))
	for i := range values {
		valuesPtrs[i] = &values[i]
	}

	var resources []Resource
	for rows.Next() {
		resource := Resource{}
		if err := rows.Scan(valuesPtrs...); err != nil {
			return nil, err
		}

		for i := range values {
			resource[queryFields[i]] = values[i]
		}

		resources = append(resources, resource)
	}
	return resources, nil
}

func (r *repository) List(ctx context.Context) ([]Collection, error) {
	rows, err := config.GetDB(ctx).QueryContext(ctx, `SELECT slug, name, fields, created_at, created_by, updated_at, updated_by FROM "collection"`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var collections []Collection
	for rows.Next() {
		var coll Collection
		if err := rows.Scan(&coll.Slug, &coll.Name, &coll.Fields, &coll.CreatedAt, &coll.CreatedBy, &coll.UpdatedAt, &coll.UpdatedBy); err != nil {
			return nil, err
		}
		collections = append(collections, coll)
	}
	return collections, nil
}
