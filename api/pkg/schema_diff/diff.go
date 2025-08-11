package schema_diff

import (
	"github.com/mimsy-cms/mimsy/pkg/schema_generator"
	"github.com/oapi-codegen/nullable"
	"github.com/xataio/pgroll/pkg/migrations"
)

func Diff(oldSchema schema_generator.SqlSchema, newSchema schema_generator.SqlSchema) []migrations.Operation {
	operations := []migrations.Operation{}

	for _, table := range newSchema.Tables {
		if oldTable, ok := oldSchema.GetTable(table.Name); ok {
			for _, column := range table.Columns {
				if oldColumn, ok := oldTable.GetColumn(column.Name); ok {
					if column.Type != oldColumn.Type {
						operation := migrations.OpAlterColumn{
							Table:  table.Name,
							Column: column.Name,
							Type:   &column.Type,
						}
						operations = append(operations, &operation)
					} else if column.IsNotNull != oldColumn.IsNotNull {
						operation := migrations.OpAlterColumn{
							Table:    table.Name,
							Column:   column.Name,
							Nullable: &column.IsNotNull,
						}
						operations = append(operations, &operation)
					} else if column.DefaultValue != oldColumn.DefaultValue {
						operation := migrations.OpAlterColumn{
							Table:   table.Name,
							Column:  column.Name,
							Default: nullable.NewNullableWithValue(column.DefaultValue),
						}
						operations = append(operations, &operation)
					}
				} else {
					operation := migrations.OpAddColumn{
						Table: table.Name,
						Column: migrations.Column{
							Name:     column.Name,
							Type:     column.Type,
							Nullable: column.IsNotNull,
						},
					}
					operations = append(operations, &operation)
				}
			}
		} else {
			columns := make([]migrations.Column, len(table.Columns))
			for i, column := range table.Columns {
				var defaultValue *string
				if column.DefaultValue != "" {
					defaultValue = &column.DefaultValue
				}

				columns[i] = migrations.Column{
					Name:     column.Name,
					Type:     column.Type,
					Nullable: column.IsNotNull,
					Default:  defaultValue,
				}
			}

			operation := migrations.OpCreateTable{
				Name:    table.Name,
				Columns: columns,
			}
			operations = append(operations, &operation)
		}
	}

	for _, oldTable := range oldSchema.Tables {
		if _, ok := newSchema.GetTable(oldTable.Name); !ok {
			operation := migrations.OpDropTable{
				Name: oldTable.Name,
			}
			operations = append(operations, &operation)
		}
	}

	for _, oldTable := range oldSchema.Tables {
		if newTable, ok := newSchema.GetTable(oldTable.Name); ok {
			for _, oldColumn := range oldTable.Columns {
				if _, ok := newTable.GetColumn(oldColumn.Name); !ok {
					operation := migrations.OpDropColumn{
						Table:  oldTable.Name,
						Column: oldColumn.Name,
					}
					operations = append(operations, &operation)
				}
			}
		}
	}

	return operations
}
