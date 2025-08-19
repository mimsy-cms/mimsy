package schema_generator

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSqlSchema_MarshalUnmarshal(t *testing.T) {
	originalSchema := &SqlSchema{
		Tables: []*Table{
			{
				Name: "users",
				Columns: []Column{
					{
						Name:         "id",
						Type:         "INTEGER",
						IsPrimaryKey: true,
						IsNotNull:    true,
					},
					{
						Name:         "email",
						Type:         "VARCHAR(255)",
						IsNotNull:    true,
						DefaultValue: "",
					},
					{
						Name:         "full_name",
						Type:         "TEXT",
						IsNotNull:    false,
						DefaultValue: "'Anonymous'",
					},
					{
						Name:        "search_vector",
						Type:        "tsvector",
						GeneratedAs: "to_tsvector('english', full_name)",
					},
				},
				Constraints: []Constraint{
					&PrimaryKeyConstraint{
						Table: "users",
						Key:   "id",
					},
					&UniqueConstraint{
						Table: "users",
						Key:   "email",
					},
				},
			},
			{
				Name: "posts",
				Columns: []Column{
					{
						Name:         "id",
						Type:         "INTEGER",
						IsPrimaryKey: true,
						IsNotNull:    true,
					},
					{
						Name:      "user_id",
						Type:      "INTEGER",
						IsNotNull: true,
					},
					{
						Name:         "title",
						Type:         "VARCHAR(500)",
						IsNotNull:    true,
						DefaultValue: "''",
					},
					{
						Name: "content",
						Type: "TEXT",
					},
				},
				Constraints: []Constraint{
					&PrimaryKeyConstraint{
						Table: "posts",
						Key:   "id",
					},
					&ForeignKeyConstraint{
						Table:           "posts",
						Column:          "user_id",
						ReferenceTable:  "users",
						ReferenceColumn: "id",
					},
				},
			},
			{
				Name: "tags",
				Columns: []Column{
					{
						Name:      "post_id",
						Type:      "INTEGER",
						IsNotNull: true,
					},
					{
						Name:      "tag_name",
						Type:      "VARCHAR(50)",
						IsNotNull: true,
					},
				},
				Constraints: []Constraint{
					&CompositePrimaryKeyConstraint{
						Table:   "tags",
						Columns: []string{"post_id", "tag_name"},
					},
					&ForeignKeyConstraint{
						Table:           "tags",
						Column:          "post_id",
						ReferenceTable:  "posts",
						ReferenceColumn: "id",
					},
				},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(originalSchema)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaledSchema SqlSchema
	err = json.Unmarshal(data, &unmarshaledSchema)
	require.NoError(t, err)

	// Verify the structure
	require.Equal(t, len(originalSchema.Tables), len(unmarshaledSchema.Tables))

	for i, originalTable := range originalSchema.Tables {
		unmarshaledTable := unmarshaledSchema.Tables[i]
		
		assert.Equal(t, originalTable.Name, unmarshaledTable.Name)
		assert.Equal(t, len(originalTable.Columns), len(unmarshaledTable.Columns))
		
		for j, originalColumn := range originalTable.Columns {
			unmarshaledColumn := unmarshaledTable.Columns[j]
			assert.Equal(t, originalColumn.Name, unmarshaledColumn.Name)
			assert.Equal(t, originalColumn.Type, unmarshaledColumn.Type)
			assert.Equal(t, originalColumn.IsPrimaryKey, unmarshaledColumn.IsPrimaryKey)
			assert.Equal(t, originalColumn.IsNotNull, unmarshaledColumn.IsNotNull)
			assert.Equal(t, originalColumn.DefaultValue, unmarshaledColumn.DefaultValue)
			assert.Equal(t, originalColumn.GeneratedAs, unmarshaledColumn.GeneratedAs)
		}
		
		// Note: Constraints will need special handling since they're interfaces
		assert.Equal(t, len(originalTable.Constraints), len(unmarshaledTable.Constraints))
	}
}

func TestSqlSchema_GetTable(t *testing.T) {
	schema := &SqlSchema{
		Tables: []*Table{
			{Name: "users"},
			{Name: "posts"},
		},
	}

	// Test existing table
	table, found := schema.GetTable("users")
	assert.True(t, found)
	assert.NotNil(t, table)
	assert.Equal(t, "users", table.Name)

	// Test non-existing table
	table, found = schema.GetTable("comments")
	assert.False(t, found)
	assert.Nil(t, table)
}

func TestTable_GetColumn(t *testing.T) {
	table := &Table{
		Name: "users",
		Columns: []Column{
			{Name: "id", Type: "INTEGER"},
			{Name: "email", Type: "VARCHAR(255)"},
		},
	}

	// Test existing column
	column, found := table.GetColumn("email")
	assert.True(t, found)
	assert.NotNil(t, column)
	assert.Equal(t, "email", column.Name)
	assert.Equal(t, "VARCHAR(255)", column.Type)

	// Test non-existing column
	column, found = table.GetColumn("password")
	assert.False(t, found)
	assert.Nil(t, column)
}

func TestSqlSchema_MarshalUnmarshal_ProducesSameSQL(t *testing.T) {
	originalSchema := &SqlSchema{
		Tables: []*Table{
			{
				Name: "users",
				Columns: []Column{
					{
						Name:         "id",
						Type:         "INTEGER",
						IsPrimaryKey: true,
						IsNotNull:    true,
					},
					{
						Name:         "email",
						Type:         "VARCHAR(255)",
						IsNotNull:    true,
					},
					{
						Name:        "search_vector",
						Type:        "tsvector",
						GeneratedAs: "to_tsvector('english', email)",
					},
				},
				Constraints: []Constraint{
					&PrimaryKeyConstraint{
						Table: "users",
						Key:   "id",
					},
					&UniqueConstraint{
						Table: "users",
						Key:   "email",
					},
				},
			},
			{
				Name: "posts",
				Columns: []Column{
					{
						Name:      "id",
						Type:      "INTEGER",
						IsNotNull: true,
					},
					{
						Name:      "user_id",
						Type:      "INTEGER",
						IsNotNull: true,
					},
					{
						Name:         "title",
						Type:         "VARCHAR(500)",
						DefaultValue: "'Untitled'",
					},
				},
				Constraints: []Constraint{
					&ForeignKeyConstraint{
						Table:           "posts",
						Column:          "user_id",
						ReferenceTable:  "mimsy_collections.users",
						ReferenceColumn: "id",
					},
				},
			},
		},
	}

	// Get original SQL
	originalSQL := originalSchema.ToSql()

	// Marshal to JSON
	data, err := json.Marshal(originalSchema)
	require.NoError(t, err)

	// Unmarshal back
	var unmarshaledSchema SqlSchema
	err = json.Unmarshal(data, &unmarshaledSchema)
	require.NoError(t, err)

	// Get SQL from unmarshaled schema
	unmarshaledSQL := unmarshaledSchema.ToSql()

	// Compare SQL outputs
	assert.Equal(t, originalSQL, unmarshaledSQL, "SQL output should be identical after marshal/unmarshal")
}