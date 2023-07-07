package builder

type sqlQuery struct {
	query string
}

type SqlQuery interface {
	Select(columns ...string) SqlQuery
	From(tables ...string) SqlQuery
	Where(conditions ...string) SqlQuery
	OrderBy(columns ...string) SqlQuery
	Join(tables ...string) SqlQuery
	On(conditions ...string) SqlQuery
	Build() string
}

func NewSqlQuery() SqlQuery {
	return &sqlQuery{
		query: "",
	}
}

func join(arr []string, sep string) string {
	var str string
	for i, v := range arr {
		if i == len(arr)-1 {
			str += v
		} else {
			str += v + sep
		}
	}
	return str
}

func (s *sqlQuery) Select(columns ...string) SqlQuery {
	s.query += "SELECT (" + join(columns, ", ") + ")"
	return s
}

func (s *sqlQuery) From(tables ...string) SqlQuery {
	s.query += " FROM " + join(tables, ", ")
	return s
}

func (s *sqlQuery) Where(conditions ...string) SqlQuery {
	s.query += " WHERE " + join(conditions, " AND ")
	return s
}

func (s *sqlQuery) OrderBy(columns ...string) SqlQuery {
	s.query += " ORDER BY " + join(columns, ", ")
	return s
}

func (s *sqlQuery) Join(tables ...string) SqlQuery {
	s.query += " JOIN " + join(tables, ", ")
	return s
}

func (s *sqlQuery) On(conditions ...string) SqlQuery {
	s.query += " ON " + join(conditions, " AND ")
	return s
}

func (s *sqlQuery) Build() string {
	return s.query
}
