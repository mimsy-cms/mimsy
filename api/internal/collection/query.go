package collection

import (
	"context"
	"fmt"
	"strings"

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

func (q *selectQuery) FindOne(ctx context.Context) (*Resource, error) {
	query := q.buildSelectQuery(q.tableName)

	row := config.GetDB(ctx).QueryRowContext(ctx, query)
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
	query := q.buildSelectQuery(q.tableName)

	rows, err := config.GetDB(ctx).QueryContext(ctx, query)
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

	resource := Resource{}
	for i := range values {
		resource[q.queryFields[i]] = values[i]
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
