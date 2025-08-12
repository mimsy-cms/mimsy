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

func TestDiffConstraintAdded(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "email", Type: "varchar(100)", IsNotNull: true},
				},
				Constraints: []schema_generator.Constraint{
					&schema_generator.PrimaryKeyConstraint{Table: "users", Key: "id"},
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
					{Name: "email", Type: "varchar(100)", IsNotNull: true},
				},
				Constraints: []schema_generator.Constraint{
					&schema_generator.PrimaryKeyConstraint{Table: "users", Key: "id"},
					&schema_generator.UniqueConstraint{Table: "users", Key: "email"},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (add constraint), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpCreateConstraint); ok {
		if op.Table != "users" || op.Type != migrations.OpCreateConstraintTypeUnique {
			t.Errorf("expected create unique constraint operation for 'email' in 'users', got %s.%s", op.Table, op.Type)
		}
	} else {
		t.Errorf("expected operation to be OpCreateConstraint, got %T", diff[0])
	}
}

func TestDiffConstraintDropped(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "users",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "email", Type: "varchar(100)", IsNotNull: true},
				},
				Constraints: []schema_generator.Constraint{
					&schema_generator.PrimaryKeyConstraint{Table: "users", Key: "id"},
					&schema_generator.UniqueConstraint{Table: "users", Key: "email"},
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
					{Name: "email", Type: "varchar(100)", IsNotNull: true},
				},
				Constraints: []schema_generator.Constraint{
					&schema_generator.PrimaryKeyConstraint{Table: "users", Key: "id"},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (drop constraint), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpDropMultiColumnConstraint); ok {
		if op.Table != "users" {
			t.Errorf("expected drop unique constraint operation for 'email' in 'users', got %s.%s", op.Table, op.Name)
		}
	} else {
		t.Errorf("expected operation to be OpDropMultiColumnConstraint, got %T", diff[0])
	}
}

func TestDiffForeignKeyConstraintAdded(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "posts",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "user_id", Type: "bigint", IsNotNull: true},
				},
				Constraints: []schema_generator.Constraint{
					&schema_generator.PrimaryKeyConstraint{Table: "posts", Key: "id"},
				},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "posts",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "user_id", Type: "bigint", IsNotNull: true},
				},
				Constraints: []schema_generator.Constraint{
					&schema_generator.PrimaryKeyConstraint{Table: "posts", Key: "id"},
					&schema_generator.ForeignKeyConstraint{
						Table:           "posts",
						Column:          "user_id",
						ReferenceTable:  "users",
						ReferenceColumn: "id",
					},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (add foreign key constraint), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpCreateConstraint); ok {
		if op.Table != "posts" || op.Type != migrations.OpCreateConstraintTypeForeignKey {
			t.Errorf("expected create foreign key constraint operation for 'user_id' in 'posts', got %s.%s", op.Table, op.Type)
		}
		if op.References == nil || op.References.Table != "users" || len(op.References.Columns) != 1 || op.References.Columns[0] != "id" {
			t.Errorf("expected foreign key to reference 'users.id', got %v", op.References)
		}
	} else {
		t.Errorf("expected operation to be OpCreateConstraint, got %T", diff[0])
	}
}

func TestDiffCompositePrimaryKeyConstraintAdded(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "user_posts",
				Columns: []schema_generator.Column{
					{Name: "user_id", Type: "bigint", IsNotNull: true},
					{Name: "post_id", Type: "bigint", IsNotNull: true},
				},
				Constraints: []schema_generator.Constraint{},
			},
		},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "user_posts",
				Columns: []schema_generator.Column{
					{Name: "user_id", Type: "bigint", IsNotNull: true},
					{Name: "post_id", Type: "bigint", IsNotNull: true},
				},
				Constraints: []schema_generator.Constraint{
					&schema_generator.CompositePrimaryKeyConstraint{
						Table:   "user_posts",
						Columns: []string{"user_id", "post_id"},
					},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (add composite primary key constraint), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpCreateConstraint); ok {
		if op.Table != "user_posts" || op.Type != migrations.OpCreateConstraintTypePrimaryKey {
			t.Errorf("expected create primary key constraint operation for 'user_posts', got %s.%s", op.Table, op.Type)
		}
		if len(op.Columns) != 2 || op.Columns[0] != "user_id" || op.Columns[1] != "post_id" {
			t.Errorf("expected primary key columns to be ['user_id', 'post_id'], got %v", op.Columns)
		}
	} else {
		t.Errorf("expected operation to be OpCreateConstraint, got %T", diff[0])
	}
}

func TestDiffCreateTableWithConstraints(t *testing.T) {
	oldSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{},
	}

	newSchema := schema_generator.SqlSchema{
		Tables: []*schema_generator.Table{
			{
				Name: "posts",
				Columns: []schema_generator.Column{
					{Name: "id", Type: "bigint", IsPrimaryKey: true, IsNotNull: true},
					{Name: "title", Type: "varchar(255)", IsNotNull: true},
					{Name: "user_id", Type: "bigint", IsNotNull: true},
				},
				Constraints: []schema_generator.Constraint{
					&schema_generator.PrimaryKeyConstraint{Table: "posts", Key: "id"},
					&schema_generator.UniqueConstraint{Table: "posts", Key: "title"},
					&schema_generator.ForeignKeyConstraint{
						Table:           "posts",
						Column:          "user_id",
						ReferenceTable:  "users",
						ReferenceColumn: "id",
					},
				},
			},
		},
	}

	diff := schema_diff.Diff(oldSchema, newSchema)
	if len(diff) != 1 {
		t.Errorf("expected 1 operation (create table), got %d", len(diff))
	}

	if op, ok := diff[0].(*migrations.OpCreateTable); ok {
		if len(op.Constraints) != 3 {
			t.Errorf("expected 3 constraints on new table, got %d", len(op.Constraints))
		}

		foundPK := false
		foundUnique := false
		foundFK := false
		for _, constraint := range op.Constraints {
			switch constraint.Type {
			case migrations.ConstraintTypePrimaryKey:
				if len(constraint.Columns) == 1 && constraint.Columns[0] == "id" {
					foundPK = true
				}
			case migrations.ConstraintTypeUnique:
				if len(constraint.Columns) == 1 && constraint.Columns[0] == "title" {
					foundUnique = true
				}
			case migrations.ConstraintTypeForeignKey:
				if len(constraint.Columns) == 1 && constraint.Columns[0] == "user_id" &&
					constraint.References != nil && constraint.References.Table == "users" {
					foundFK = true
				}
			}
		}

		if !foundPK {
			t.Errorf("expected primary key constraint not found")
		}
		if !foundUnique {
			t.Errorf("expected unique constraint not found")
		}
		if !foundFK {
			t.Errorf("expected foreign key constraint not found")
		}
	} else {
		t.Errorf("expected operation to be OpCreateTable, got %T", diff[0])
	}
}
