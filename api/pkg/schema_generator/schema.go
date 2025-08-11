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

func (s *SqlSchema) GetTable(name string) (*Table, bool) {
	for _, table := range s.Tables {
		if table.Name == name {
			return table, true
		}
	}
	return nil, false
}
