package collection

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/mimsy-cms/mimsy/internal/config"
)

type rowScanner interface {
	Scan(dest ...any) error
}

type selectQuery struct {
	tableName   string
	queryFields []string
}

func NewSelectQuery(tableName string, fields map[string]Field) *selectQuery {
	return &selectQuery{
		tableName:   tableName,
		queryFields: transformQueryFields(fields),
	}
}

var (
	// defaultColumns are the columns that will always be selected in a query.
	defaultColumns = []string{"id", "slug"}
)

func (q *selectQuery) FindOne(ctx context.Context, slug string) (*Resource, error) {
	query := fmt.Sprintf(
		`SELECT id, slug, content, created_at, updated_at FROM %s WHERE slug = $1`,
		pq.QuoteIdentifier(q.tableName),
	)

	var resource Resource
	err := config.GetDB(ctx).QueryRowContext(ctx, query, slug).Scan(
		&resource.ID,
		&resource.Slug,
		&resource.Content,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return &resource, nil
}

func (q *selectQuery) FindAll(ctx context.Context) ([]Resource, error) {
	query := fmt.Sprintf(
		`SELECT id, slug, content, created_at, updated_at FROM %s ORDER BY created_at DESC`,
		pq.QuoteIdentifier(q.tableName),
	)

	rows, err := config.GetDB(ctx).QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var resources []Resource
	for rows.Next() {
		var resource Resource
		if err := rows.Scan(
			&resource.ID,
			&resource.Slug,
			&resource.Content,
			&resource.CreatedAt,
			&resource.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		resources = append(resources, resource)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return resources, nil
}

func (q *selectQuery) scan(row rowScanner) (*Resource, error) {
	values := make([]any, len(q.queryFields))
	valuesPtrs := make([]any, len(q.queryFields))
	for i := range values {
		valuesPtrs[i] = &values[i]
	}

	if err := row.Scan(valuesPtrs...); err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	resource := Resource{}
	for i := range values {
		switch q.queryFields[i] {
		case "id":
			if v, ok := values[i].(int64); ok {
				resource.ID = int(v)
			}
		case "slug":
			if v, ok := values[i].(string); ok {
				resource.Slug = v
			}
		case "content":
			if v, ok := values[i].(string); ok {
				resource.Content = json.RawMessage(v)
			}
		case "created_at":
			if v, ok := values[i].(time.Time); ok {
				resource.CreatedAt = v
			}
		case "updated_at":
			if v, ok := values[i].(time.Time); ok {
				resource.UpdatedAt = v
			}
			// Add more fields as needed
		}
	}

	return &resource, nil
}

func transformQueryFields(fields map[string]Field) []string {
	queryFields := defaultColumns
	for name, field := range fields {
		if field.Type == "relation" {
			if field.Relation != nil && field.Relation.Type == FieldRelationTypeManyToOne {
				queryFields = append(queryFields, fmt.Sprintf("%s_id", name))
			}
		} else {
			queryFields = append(queryFields, name)
		}
	}
	return queryFields
}

func (q *selectQuery) buildSelectQuery(tableName string) string {
	quotedQueryFields := make([]string, len(q.queryFields))
	for i, field := range q.queryFields {
		quotedQueryFields[i] = pq.QuoteIdentifier(field)
	}

	query := fmt.Sprintf(
		`SELECT %s FROM %s`,
		strings.Join(quotedQueryFields, ", "),
		pq.QuoteIdentifier(tableName),
	)

	return query
}
