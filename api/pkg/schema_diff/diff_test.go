package schema_diff_test

import (
	"testing"

	"github.com/mimsy-cms/mimsy/pkg/schema_diff"
	"github.com/mimsy-cms/mimsy/pkg/schema_generator"
	"github.com/xataio/pgroll/pkg/migrations"
)

func TestDiff(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
					{Name: "email", Type: "varchar(100)", IsNotNull: true},
				},
			},
			{
				Name: "orders",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "user_id", Type: "bigint", IsNotNull: true},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 2 {
		t.Errorf("expected 2 operations, got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpAddColumn); ok {
		if op.Table != "users" || op.Column.Name != "email" {
			t.Errorf("expected add column operation for 'email' in 'users', got %s.%s", op.Table, op.Column.Name)
		}
	} else {
		t.Errorf("expected operation to be OpAddColumn, got %T", diff[0])
	}
}

func TestDiffRemovedTable(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (drop table), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpDropTable); ok {
		if op.Name != "users" {
			t.Errorf("expected drop table operation for 'users', got %s", op.Name)
		}
	} else {
		t.Errorf("expected operation to be OpDropTable, got %T", diff[0])
	}
}

func TestDiffColumnTypeChange(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "text", IsNotNull: true},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (alter column), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpAlterColumn); ok {
		if op.Table != "users" || op.Column != "name" || *op.Type != "text" {
			t.Errorf("expected alter column operation for 'name' in 'users' to change type to 'text', got %s.%s with type %s", op.Table, op.Column, *op.Type)
		}
	} else {
		t.Errorf("expected operation to be OpAlterColumn, got %T", diff[0])
	}
}

func TestDiffColumnAdded(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
					{Name: "email", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (add column), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpAddColumn); ok {
		if op.Table != "users" || op.Column.Name != "email" {
			t.Errorf("expected add column operation for 'email' in 'users', got %s.%s", op.Table, op.Column.Name)
		}
	} else {
		t.Errorf("expected operation to be OpAddColumn, got %T", diff[0])
	}
}

func TestDiffColumnRemoved(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
					{Name: "email", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (drop column), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpDropColumn); ok {
		if op.Table != "users" || op.Column != "email" {
			t.Errorf("expected drop column operation for 'email' in 'users', got %s.%s", op.Table, op.Column)
		}
	} else {
		t.Errorf("expected operation to be OpDropColumn, got %T", diff[0])
	}
}

func TestDiffNoChanges(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 0 {
		t.Errorf("expected no operations, got %d", len(diff))
	}
}

func TestDefaultValueChange(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true, DefaultValue: "John Doe"},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true, DefaultValue: "Jane Doe"},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (alter column default value), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpAlterColumn); ok {
		if op.Default == nil {
			t.Errorf("expected default value to be 'Jane Doe', got %v", op.Default)
		}
	} else {
		t.Errorf("expected operation to be OpAlterColumn, got %T", diff[0])
	}

	if op, ok := diff[0].(*migrations.OpAlterColumn); ok {
		if op.Default == nil {
			t.Errorf("expected default value to be 'Jane Doe', got %v", op.Default)
		}
	} else {
		t.Errorf("expected operation to be OpAlterColumn, got %T", diff[0])
	}
}

func TestDiffAddDefaultValue(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true, DefaultValue: "Default Name"},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (add default value), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpAlterColumn); ok {
		if op.Default == nil {
			t.Errorf("expected default value to be 'Default Name', got %v", op.Default)
		}
		value, err := op.Default.Get()
		if err != nil || value != "Default Name" {
			t.Errorf("expected default value to be 'Default Name', got %v", value)
		}
	} else {
		t.Errorf("expected operation to be OpAlterColumn, got %T", diff[0])
	}
}

func TestDiffRemoveDefaultValue(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true, DefaultValue: "Default Name"},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (remove default value), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpAlterColumn); ok {
		if op.Default == nil {
			t.Errorf("expected default value to be removed, got %v", op.Default)
		}
	} else {
		t.Errorf("expected operation to be OpAlterColumn, got %T", diff[0])
	}
}

func TestDiffColumnNullableChange(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: true},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "name", Type: "varchar(100)", IsNotNull: false},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (alter column nullable), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpAlterColumn); ok {
		if op.Table != "users" || op.Column != "name" || op.Nullable == nil || *op.Nullable != false {
			t.Errorf("expected alter column operation for 'name' in 'users' to change nullable to false, got %s.%s with nullable %v", op.Table, op.Column, *op.Nullable)
		}
	} else {
		t.Errorf("expected operation to be OpAlterColumn, got %T", diff[0])
	}
}
