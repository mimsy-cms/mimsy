package schema_diff

import (
	"log/slog"
	"strings"

	"github.com/mimsy-cms/mimsy/pkg/schema_generator"
	"github.com/oapi-codegen/nullable"
	"github.com/xataio/pgroll/pkg/migrations"
)

func Diff(oldSchema schema_generator.SqlSchema, newSchema schema_generator.SqlSchema) []migrations.Operation {
	operations := []migrations.Operation{}

	operations = append(operations, processTableChanges(oldSchema, newSchema)...)
	operations = append(operations, processConstraintChanges(oldSchema, newSchema)...)
	operations = append(operations, processDroppedTables(oldSchema, newSchema)...)
	operations = append(operations, processDroppedColumns(oldSchema, newSchema)...)
	operations = append(operations, processDroppedConstraints(oldSchema, newSchema)...)

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

		var generated *migrations.ColumnGenerated
		if column.IsPrimaryKey {
			generated = &migrations.ColumnGenerated{
				Identity: &migrations.ColumnGeneratedIdentity{UserSpecifiedValues: "BY DEFAULT"},
			}
		}

		columns[i] = migrations.Column{
			Name:      column.Name,
			Type:      column.Type,
			Nullable:  column.IsNotNull,
			Default:   defaultValue,
			Generated: generated,
		}
	}

	constraints := make([]migrations.Constraint, 0, len(table.Constraints))
	for _, constraint := range table.Constraints {
		tableConstraint := createTableConstraint(constraint)
		if tableConstraint != nil {
			constraints = append(constraints, *tableConstraint)
		}
	}

	return &migrations.OpCreateTable{
		Name:        table.Name,
		Columns:     columns,
		Constraints: constraints,
	}
}

func createTableConstraint(constraint schema_generator.Constraint) *migrations.Constraint {
	switch c := constraint.(type) {
	case *schema_generator.UniqueConstraint:
		return &migrations.Constraint{
			Name:    c.Name(),
			Type:    migrations.ConstraintTypeUnique,
			Columns: []string{c.Key},
		}
	case *schema_generator.PrimaryKeyConstraint:
		return &migrations.Constraint{
			Name:    c.Name(),
			Type:    migrations.ConstraintTypePrimaryKey,
			Columns: []string{c.Key},
		}
	case *schema_generator.CompositePrimaryKeyConstraint:
		return &migrations.Constraint{
			Name:    c.Name(),
			Type:    migrations.ConstraintTypePrimaryKey,
			Columns: c.Columns,
		}
	case *schema_generator.ForeignKeyConstraint:
		// For cross-schema references, pgroll needs just the table name
		// when the schema is in the search path
		refTable := strings.ReplaceAll(c.ReferenceTable, "\"", "")
		// Remove schema prefix if it's mimsy_internal or mimsy_collections (which are in search path)
		if strings.HasPrefix(refTable, "mimsy_internal.") {
			refTable = strings.TrimPrefix(refTable, "mimsy_internal.")
		} else if strings.HasPrefix(refTable, "mimsy_collections.") {
			refTable = strings.TrimPrefix(refTable, "mimsy_collections.")
		}
		return &migrations.Constraint{
			Name:    c.Name(),
			Type:    migrations.ConstraintTypeForeignKey,
			Columns: []string{c.Column},
			References: &migrations.TableForeignKeyReference{
				Table:   refTable,
				Columns: []string{c.ReferenceColumn},
			},
		}
	default:
		return nil
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

func processConstraintChanges(oldSchema, newSchema schema_generator.SqlSchema) []migrations.Operation {
	operations := []migrations.Operation{}

	for _, table := range newSchema.Tables {
		oldTable, exists := oldSchema.GetTable(table.Name)
		if !exists {
			continue
		}

		operations = append(operations, processTableConstraintChanges(table, oldTable)...)
	}

	return operations
}

func processTableConstraintChanges(table, oldTable *schema_generator.Table) []migrations.Operation {
	operations := []migrations.Operation{}

	for _, constraint := range table.Constraints {
		if !constraintExists(oldTable.Constraints, constraint) {
			operation := createConstraintOperation(table.Name, constraint)
			if operation != nil {
				operations = append(operations, operation)
			}
		}
	}

	return operations
}

func constraintExists(constraints []schema_generator.Constraint, targetConstraint schema_generator.Constraint) bool {
	for _, constraint := range constraints {
		if constraint.Name() == targetConstraint.Name() {
			return true
		}
	}
	return false
}

func createConstraintOperation(tableName string, constraint schema_generator.Constraint) migrations.Operation {
	switch c := constraint.(type) {
	case *schema_generator.UniqueConstraint:
		return &migrations.OpCreateConstraint{
			Type:    migrations.OpCreateConstraintTypeUnique,
			Name:    c.Name(),
			Table:   tableName,
			Columns: []string{c.Key},
			Up: map[string]string{
				c.Key: c.Key,
			},
			Down: map[string]string{
				c.Key: c.Key,
			},
		}
	case *schema_generator.PrimaryKeyConstraint:
		return &migrations.OpCreateConstraint{
			Type:    migrations.OpCreateConstraintTypePrimaryKey,
			Name:    c.Name(),
			Table:   tableName,
			Columns: []string{c.Key},
			Up: map[string]string{
				c.Key: c.Key,
			},
			Down: map[string]string{
				c.Key: c.Key,
			},
		}
	case *schema_generator.CompositePrimaryKeyConstraint:
		upDown := make(map[string]string)
		for _, col := range c.Columns {
			upDown[col] = col
		}
		return &migrations.OpCreateConstraint{
			Type:    migrations.OpCreateConstraintTypePrimaryKey,
			Name:    c.Name(),
			Table:   tableName,
			Columns: c.Columns,
			Up:      upDown,
			Down:    upDown,
		}
	case *schema_generator.ForeignKeyConstraint:
		// For cross-schema references, pgroll needs just the table name
		// when the schema is in the search path
		refTable := strings.ReplaceAll(c.ReferenceTable, "\"", "")
		// Remove schema prefix if it's mimsy_internal or mimsy_collections (which are in search path)
		if strings.HasPrefix(refTable, "mimsy_internal.") {
			refTable = strings.TrimPrefix(refTable, "mimsy_internal.")
		} else if strings.HasPrefix(refTable, "mimsy_collections.") {
			refTable = strings.TrimPrefix(refTable, "mimsy_collections.")
		}
		slog.Info("Foreign key", "originalReferenceTable", c.ReferenceTable, "adjustedReferenceTable", refTable)
		return &migrations.OpCreateConstraint{
			Type:    migrations.OpCreateConstraintTypeForeignKey,
			Name:    c.Name(),
			Table:   tableName,
			Columns: []string{c.Column},
			References: &migrations.TableForeignKeyReference{
				Table:   refTable,
				Columns: []string{c.ReferenceColumn},
			},
			Up: map[string]string{
				c.Column: c.Column,
			},
			Down: map[string]string{
				c.Column: c.Column,
			},
		}
	default:
		return nil
	}
}

func processDroppedConstraints(oldSchema, newSchema schema_generator.SqlSchema) []migrations.Operation {
	operations := []migrations.Operation{}

	for _, oldTable := range oldSchema.Tables {
		newTable, exists := newSchema.GetTable(oldTable.Name)
		if !exists {
			// Skip dropped tables as their constraints are handled in processDroppedTables
			continue
		}

		for _, constraint := range oldTable.Constraints {
			if !constraintExists(newTable.Constraints, constraint) {
				operation := createDropConstraintOperation(oldTable.Name, constraint)
				if operation != nil {
					operations = append(operations, operation)
				}
			}
		}
	}

	return operations
}

func createDropConstraintOperation(tableName string, constraint schema_generator.Constraint) migrations.Operation {
	switch c := constraint.(type) {
	case *schema_generator.UniqueConstraint:
		return &migrations.OpDropMultiColumnConstraint{
			Name:  c.Name(),
			Table: tableName,
			Up: map[string]string{
				c.Key: c.Key,
			},
			Down: map[string]string{
				c.Key: c.Key,
			},
		}
	case *schema_generator.ForeignKeyConstraint:
		return &migrations.OpDropMultiColumnConstraint{
			Name:  c.Name(),
			Table: tableName,
			Up: map[string]string{
				c.Column: c.Column,
			},
			Down: map[string]string{
				c.Column: c.Column,
			},
		}
	default:
		// Primary key constraints cannot be dropped easily, skip for now
		return nil
	}
}
