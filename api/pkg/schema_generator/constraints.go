package schema_generator

import (
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type Constraint interface {
	Name() string
	ToSql() string
}

type UniqueConstraint struct {
	Table string
	Key   string
}

func (u *UniqueConstraint) Name() string {
	return fmt.Sprintf("uq__%s__%s", u.Table, u.Key)
}

func (u *UniqueConstraint) ToSql() string {
	return fmt.Sprintf("CONSTRAINT %s UNIQUE (%s)", u.Name(), pq.QuoteIdentifier(u.Key))
}

type PrimaryKeyConstraint struct {
	Table string
	Key   string
}

func (p *PrimaryKeyConstraint) Name() string {
	return fmt.Sprintf("pk__%s", p.Table)
}

func (p *PrimaryKeyConstraint) ToSql() string {
	return fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", p.Name(), pq.QuoteIdentifier(p.Key))
}

type CompositePrimaryKeyConstraint struct {
	Table   string
	Columns []string
}

func (c *CompositePrimaryKeyConstraint) Name() string {
	return fmt.Sprintf("pk__%s", c.Table)
}

func (c *CompositePrimaryKeyConstraint) ToSql() string {
	quotedColumns := make([]string, len(c.Columns))
	for i, column := range c.Columns {
		quotedColumns[i] = pq.QuoteIdentifier(column)
	}

	return fmt.Sprintf("CONSTRAINT %s PRIMARY KEY (%s)", c.Name(), strings.Join(quotedColumns, ", "))
}

type ForeignKeyConstraint struct {
	Table           string
	Column          string
	ReferenceColumn string
	ReferenceTable  string
}

func (f *ForeignKeyConstraint) Name() string {
	return fmt.Sprintf("fk__%s__%s__%s", f.Table, f.Column, f.ReferenceTable)
}

func (f *ForeignKeyConstraint) ToSql() string {
	return fmt.Sprintf("CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s (%s)",
		f.Name(),
		pq.QuoteIdentifier(f.Column),
		pq.QuoteIdentifier(f.ReferenceTable),
		pq.QuoteIdentifier(f.ReferenceColumn),
	)
}
