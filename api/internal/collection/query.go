package collection

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
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
	defaultColumns  = []string{"id", "slug", "created_at", "created_by", "updated_at", "updated_by"}
	readOnlyColumns = []string{"id", "created_at", "created_by", "updated_at", "updated_by"}
)

func (q *selectQuery) FindOne(ctx context.Context, slug string) (*Resource, error) {
	query, args, err := q.buildSelectQuery(q.tableName, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to build select SQL query: %w", err)
	}

	row := config.GetDB(ctx).QueryRowContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, fmt.Errorf("failed to execute query: %w", row.Err())
	}

	resource, err := q.scan(row)
	if err != nil {
		return nil, err
	}
	return resource, nil
}

func (q *selectQuery) FindAll(ctx context.Context) ([]Resource, error) {
	query, args, err := q.buildSelectQuery(q.tableName)
	if err != nil {
		return nil, fmt.Errorf("failed to build select SQL query: %w", err)
	}

	rows, err := config.GetDB(ctx).QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	resources := []Resource{}
	for rows.Next() {
		resource, err := q.scan(rows)
		if err != nil {
			return nil, err
		}
		resources = append(resources, *resource)
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

	resource := Resource{Fields: map[string]any{}}
	for i := 4; i < len(values); i++ {
		resource.Fields[q.queryFields[i]] = values[i]
	}

	resource.Id = values[0].(int64)
	resource.Slug = values[1].(string)
	resource.CreatedAt = values[2].(time.Time)
	resource.CreatedBy = values[3].(string)
	resource.UpdatedAt = values[4].(time.Time)
	resource.UpdatedBy = values[5].(string)

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

func (q *selectQuery) buildSelectQuery(tableName string, slug ...string) (string, []any, error) {
	quotedQueryFields := make([]string, len(q.queryFields))
	for i, field := range q.queryFields {
		quotedQueryFields[i] = pq.QuoteIdentifier(field)
	}

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	b := psql.Select(quotedQueryFields...).From(pq.QuoteIdentifier(tableName))

	if len(slug) > 0 {
		b = b.Where(sq.Eq{"slug": slug[0]})
	}

	return b.ToSql()
}
