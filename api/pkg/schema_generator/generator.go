package schema_generator

import (
	"fmt"
	"slices"

	"github.com/lib/pq"
	"github.com/mimsy-cms/mimsy/pkg/mimsy_schema"
)

type SchemaGenerator interface {
	GenerateSqlSchema(*mimsy_schema.Schema) (SqlSchema, error)
}

type schemaGenerator struct {
}

func New() SchemaGenerator {
	return &schemaGenerator{}
}
func (s *schemaGenerator) GenerateSqlSchema(schema *mimsy_schema.Schema) (SqlSchema, error) {
	// For each collection in the schema
	sqlSchema := SqlSchema{}
	for _, collection := range schema.Collections {
		collectionSqlSchema, err := s.HandleCollection(collection)
		if err != nil {
			return SqlSchema{}, err
		}

		sqlSchema = MergeSchemas(sqlSchema, &collectionSqlSchema)
	}

	return sqlSchema, nil
}

func (s *schemaGenerator) GenerateIdColumn() Column {
	return Column{
		Name:         "id",
		Type:         "bigint",
		IsPrimaryKey: true,
		IsNotNull:    true,
	}
}

func (s *schemaGenerator) GenerateIdConstraint(collection *mimsy_schema.Collection) Constraint {
	return &PrimaryKeyConstraint{
		Table: collection.Name,
		Key:   "id",
	}
}

func (s *schemaGenerator) GenerateSlugColumn() Column {
	return Column{
		Name:      "slug",
		Type:      "varchar(60)",
		IsNotNull: true,
	}
}

func (s *schemaGenerator) GenerateSlugConstraint(collection *mimsy_schema.Collection) Constraint {
	return &UniqueConstraint{
		Table: collection.Name,
		Key:   "slug",
	}
}

type Entry struct {
	Name  string
	Value mimsy_schema.SchemaElement
}

func OrderFields(collection *mimsy_schema.Collection) []Entry {
	var fields []Entry

	for name, element := range collection.Schema {
		fields = append(fields, Entry{Name: name, Value: element})
	}

	slices.SortFunc(fields, func(a, b Entry) int {
		// Always put relations at the end
		if a.Value.IsRelation() && !b.Value.IsRelation() {
			return 1
		}
		if !a.Value.IsRelation() && b.Value.IsRelation() {
			return -1
		}
		if a.Name < b.Name {
			return -1
		}
		if a.Name > b.Name {
			return 1
		}
		return 0
	})

	return fields
}

func (s *schemaGenerator) HandleCollection(collection mimsy_schema.Collection) (SqlSchema, error) {
	baseTable := Table{
		Name:        collection.Name,
		Columns:     []Column{},
		Constraints: []Constraint{},
	}

	baseTable.Columns = append(baseTable.Columns,
		s.GenerateIdColumn(),
		s.GenerateSlugColumn(),
	)

	baseTable.Constraints = append(baseTable.Constraints,
		s.GenerateIdConstraint(&collection),
		s.GenerateSlugConstraint(&collection),
	)

	schema := SqlSchema{
		Tables: []*Table{&baseTable},
	}
	for _, entry := range OrderFields(&collection) {
		name, element := entry.Name, entry.Value

		if !element.IsRelation() {
			// So it is a direct field
			column, err := s.HandleDirectField(name, element)
			if err != nil {
				return SqlSchema{}, err
			}
			baseTable.Columns = append(baseTable.Columns, column)

			continue
		}

		// Handle a relation field
		relationSchema, err := s.HandleRelationField(name, element, &baseTable)
		if err != nil {
			return SqlSchema{}, err
		}

		schema = MergeSchemas(schema, relationSchema)
	}

	return schema, nil
}

func (s *schemaGenerator) HandleDirectField(name string, element mimsy_schema.SchemaElement) (Column, error) {
	switch element.Type {
	case "string":
		return Column{
			Name:         name,
			Type:         "varchar",
			IsNotNull:    element.IsRequired(),
			DefaultValue: "",
		}, nil
	case "rich_text":
		return Column{
			Name:         name,
			Type:         "jsonb",
			IsNotNull:    element.IsRequired(),
			DefaultValue: "",
		}, nil
	case "created_at":
		return Column{
			Name:         name,
			Type:         "timestamptz",
			IsNotNull:    element.IsRequired(),
			DefaultValue: "CURRENT_TIMESTAMP",
		}, nil
	default:
		return Column{}, fmt.Errorf("unsupported type: %s", element.Type)
	}
}

func (s *schemaGenerator) HandleRelationField(name string, element mimsy_schema.SchemaElement, table *Table) (*SqlSchema, error) {
	switch element.Type {
	case "relation":
		return s.HandleManyToOneField(name, element, table)
	case "multi_relation":
		return s.HandleManyToManyField(name, element, table)
	default:
		return nil, fmt.Errorf("unsupported type: %s", element.Type)
	}
}

func (s *schemaGenerator) HandleManyToManyField(name string, element mimsy_schema.SchemaElement, table *Table) (*SqlSchema, error) {
	joinTableName, err := GetRelationTableName(&element, table.Name, name)
	if err != nil {
		return nil, err
	}

	referenceTableName, err := GetSimpleTableName(element.RelatesTo)

	joinTableIdentifier := joinTableName
	idColumnName := fmt.Sprintf("%s_id", referenceTableName)
	referenceTable, err := GetPrefixedTableName(element.RelatesTo)

	if err != nil {
		return nil, err
	}

	rowId := fmt.Sprintf("%s_id", table.Name)
	relatesToId := fmt.Sprintf("%s_id", referenceTableName)
	relatesToSlug := fmt.Sprintf("%s_slug", referenceTableName)

	baseTableName, err := GetPrefixedTableName(table.Name)
	if err != nil {
		return nil, err
	}

	joinTable := Table{
		Name: joinTableIdentifier,
		Columns: []Column{
			{
				Name:      rowId,
				Type:      "bigint",
				IsNotNull: true,
			},
			{
				Name:      relatesToId,
				Type:      "bigint",
				IsNotNull: true,
			},
			{
				Name:        relatesToSlug,
				Type:        "varchar",
				IsNotNull:   true,
				GeneratedAs: fmt.Sprintf("SELECT slug FROM %s WHERE id = %s", referenceTable, pq.QuoteIdentifier(idColumnName)),
			},
		},
		Constraints: []Constraint{
			&CompositePrimaryKeyConstraint{
				Table:   joinTableIdentifier,
				Columns: []string{rowId, relatesToId},
			},
			&ForeignKeyConstraint{
				Table:           joinTableIdentifier,
				Column:          rowId,
				ReferenceTable:  baseTableName,
				ReferenceColumn: "id",
			},
			&ForeignKeyConstraint{
				Table:           joinTableIdentifier,
				Column:          relatesToId,
				ReferenceTable:  referenceTable,
				ReferenceColumn: "id",
			},
		},
	}

	return &SqlSchema{
		Tables: []*Table{
			&joinTable,
		},
	}, nil
}

func (s *schemaGenerator) HandleManyToOneField(name string, element mimsy_schema.SchemaElement, table *Table) (*SqlSchema, error) {
	// Add a new column and constraint
	idColumnName := fmt.Sprintf("%s_id", name)
	referenceTable, err := GetPrefixedTableName(element.RelatesTo)

	if err != nil {
		return nil, err
	}

	idColumn := Column{
		Name:      idColumnName,
		Type:      "bigint",
		IsNotNull: element.IsRequired(),
	}

	slugColumn := Column{
		Name:        fmt.Sprintf("%s_slug", name),
		Type:        "varchar",
		GeneratedAs: fmt.Sprintf("SELECT slug FROM %s WHERE id = %s", referenceTable, pq.QuoteIdentifier(idColumnName)),
	}

	table.Columns = append(table.Columns, idColumn, slugColumn)

	foreignKeyConstraint := &ForeignKeyConstraint{
		Table:           table.Name,
		Column:          idColumnName,
		ReferenceTable:  referenceTable,
		ReferenceColumn: "id",
	}

	table.Constraints = append(table.Constraints, foreignKeyConstraint)

	return nil, nil
}

func MergeSchemas(mainSchema SqlSchema, schemas ...*SqlSchema) SqlSchema {
	for _, schema := range schemas {
		if schema == nil {
			continue
		}

		mainSchema.Tables = append(mainSchema.Tables, schema.Tables...)
	}

	return mainSchema
}
