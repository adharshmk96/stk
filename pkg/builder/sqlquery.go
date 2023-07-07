package builder

import (
	"strings"
)

type SqlQuery interface {
	InsertInto(table string) SqlQuery
	Update(table string) SqlQuery
	DeleteFrom(table string) SqlQuery
	Fields(fields ...string) SqlQuery
	Values(values ...string) SqlQuery
	Set(values ...string) SqlQuery

	Select(columns ...string) SqlQuery
	From(tables ...string) SqlQuery
	Where(conditions ...string) SqlQuery
	OrderBy(columns ...string) SqlQuery
	Join(tables ...string) SqlQuery
	On(conditions ...string) SqlQuery
	Build() string
}

type sqlQuery struct {
	query strings.Builder
	parts []string
}

func NewSqlQuery() SqlQuery {
	return &sqlQuery{}
}

func (b *sqlQuery) InsertInto(table string) SqlQuery {
	part := "INSERT INTO " + table
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) Fields(fields ...string) SqlQuery {
	part := "(" + strings.Join(fields, ", ") + ")"
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) Update(table string) SqlQuery {
	part := "UPDATE " + table
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) DeleteFrom(table string) SqlQuery {
	part := "DELETE FROM " + table
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) Set(values ...string) SqlQuery {
	part := "SET " + strings.Join(values, ", ")
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) Values(values ...string) SqlQuery {
	part := "VALUES (" + strings.Join(values, ", ") + ")"
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) Select(columns ...string) SqlQuery {
	part := "SELECT " + strings.Join(columns, ", ")
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) From(tables ...string) SqlQuery {
	part := "FROM " + strings.Join(tables, ", ")
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) Where(conditions ...string) SqlQuery {
	part := "WHERE " + strings.Join(conditions, " AND ")
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) OrderBy(columns ...string) SqlQuery {
	part := "ORDER BY " + strings.Join(columns, ", ")
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) Join(tables ...string) SqlQuery {
	part := "JOIN " + strings.Join(tables, " JOIN ")
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) On(conditions ...string) SqlQuery {
	part := "ON " + strings.Join(conditions, " AND ")
	b.parts = append(b.parts, part)
	return b
}

func (b *sqlQuery) Build() string {
	query := strings.Join(b.parts, " ")
	b.query.WriteString(query)
	queryString := b.query.String()
	b.query.Reset()
	return queryString
}
