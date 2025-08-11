package schema_generator

type SqlSchema struct {
	Tables []*Table
}

func (s *SqlSchema) ToSql() string {
	sql := ""

	for _, table := range s.Tables {
		sql += table.ToSql()
		sql += "\n"
	}

	return sql
}
