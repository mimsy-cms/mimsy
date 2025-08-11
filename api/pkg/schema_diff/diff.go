package schema_diff

import (
	"github.com/mimsy-cms/mimsy/pkg/schema_generator"
	"github.com/oapi-codegen/nullable"
	"github.com/xataio/pgroll/pkg/migrations"
)

func Diff(oldSchema schema_generator.SqlSchema, newSchema schema_generator.SqlSchema) []migrations.Operation {
	operations := []migrations.Operation{}

	operations = append(operations, processTableChanges(oldSchema, newSchema)...)
	operations = append(operations, processDroppedTables(oldSchema, newSchema)...)
	operations = append(operations, processDroppedColumns(oldSchema, newSchema)...)

	return operations
}

func processTableChanges(oldSchema, newSchema schema_generator.SqlSchema) []migrations.Operation {
	operations := []migrations.Operation{}

	for _, table := range newSchema.Tables {
		oldTable, exists := oldSchema.GetTable(table.Name)
		if !exists {
			operations = append(operations, createTableOperation(table))
			continue
		}

		operations = append(operations, processColumnChanges(table, oldTable)...)
	}

	return operations
}

func createTableOperation(table *schema_generator.Table) migrations.Operation {
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

	return &migrations.OpCreateTable{
		Name:    table.Name,
		Columns: columns,
	}
}

func processColumnChanges(table, oldTable *schema_generator.Table) []migrations.Operation {
	operations := []migrations.Operation{}

	for _, column := range table.Columns {
		oldColumn, exists := oldTable.GetColumn(column.Name)
		if !exists {
			operation := migrations.OpAddColumn{
				Table: table.Name,
				Column: migrations.Column{
					Name:     column.Name,
					Type:     column.Type,
					Nullable: column.IsNotNull,
				},
			}
			operations = append(operations, &operation)
			continue
		}

		if alterOp := createAlterColumnOperation(table.Name, column, oldColumn); alterOp != nil {
			operations = append(operations, alterOp)
		}
	}

	return operations
}

func createAlterColumnOperation(tableName string, column schema_generator.Column, oldColumn *schema_generator.Column) migrations.Operation {
	if column.Type != oldColumn.Type {
		return &migrations.OpAlterColumn{
			Table:  tableName,
			Column: column.Name,
			Type:   &column.Type,
		}
	}

	if column.IsNotNull != oldColumn.IsNotNull {
		return &migrations.OpAlterColumn{
			Table:    tableName,
			Column:   column.Name,
			Nullable: &column.IsNotNull,
		}
	}

	if column.DefaultValue != oldColumn.DefaultValue {
		return &migrations.OpAlterColumn{
			Table:   tableName,
			Column:  column.Name,
			Default: nullable.NewNullableWithValue(column.DefaultValue),
		}
	}

	return nil
}

func processDroppedTables(oldSchema, newSchema schema_generator.SqlSchema) []migrations.Operation {
	operations := []migrations.Operation{}

	for _, oldTable := range oldSchema.Tables {
		if _, exists := newSchema.GetTable(oldTable.Name); !exists {
			operation := migrations.OpDropTable{
				Name: oldTable.Name,
			}
			operations = append(operations, &operation)
		}
	}

	return operations
}

func processDroppedColumns(oldSchema, newSchema schema_generator.SqlSchema) []migrations.Operation {
	operations := []migrations.Operation{}

	for _, oldTable := range oldSchema.Tables {
		newTable, exists := newSchema.GetTable(oldTable.Name)
		if !exists {
			continue
		}

		for _, oldColumn := range oldTable.Columns {
			if _, exists := newTable.GetColumn(oldColumn.Name); !exists {
				operation := migrations.OpDropColumn{
					Table:  oldTable.Name,
					Column: oldColumn.Name,
				}
				operations = append(operations, &operation)
			}
		}
	}

	return operations
}
