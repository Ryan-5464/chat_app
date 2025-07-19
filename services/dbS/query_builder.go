package dbService

import (
	"strings"
)

func NewQueryBuilder() *queryBuilder {
	qb := &queryBuilder{
		b: strings.Builder{},
	}
	return qb
}

type queryBuilder struct {
	b strings.Builder
}

type SQLQuery string

// type field string
type on string
type placeholders string
type values string

type field struct {
	name string
}

func (f *field) String() string {
	return f.name
}

type table struct {
	name string
}

func (t *table) String() string {
	return t.name
}

func (q *queryBuilder) Build() SQLQuery {
	s := strings.TrimSpace(q.b.String())
	return SQLQuery(s)
}

// =============================================================================
// Objects

func (q *queryBuilder) Field(name string) field {
	return field{name: name}
}

func (q *queryBuilder) FieldWithAlias(name string, tag string) field {
	var b strings.Builder
	b.WriteString(tag)
	b.WriteString(".")
	b.WriteString(name)
	return field{name: b.String()}
}

func (q *queryBuilder) All() field {
	return field{name: " * "}
}

func (q *queryBuilder) Table(name string) table {
	return table{name: name}
}

func (q *queryBuilder) TableWithAlias(name string, tag string) table {
	var b strings.Builder
	b.WriteString(name)
	b.WriteString(" AS ")
	b.WriteString(tag)
	return table{name: b.String()}
}

// =============================================================================
// Functions

func (q *queryBuilder) Count(f field) field {
	var b strings.Builder
	b.WriteString(" COUNT(")
	b.WriteString(f.String())
	b.WriteString(") AS count")
	return field{name: b.String()}
}

func (q *queryBuilder) Max(f field) field {
	var b strings.Builder
	b.WriteString(" MAX(")
	b.WriteString(f.String())
	b.WriteString(") AS max")
	return field{name: b.String()}
}

// =============================================================================
// Conditions

func (q *queryBuilder) EqualTo() placeholders {
	return placeholders("= ?")
}

func (q *queryBuilder) GreaterThan() placeholders {
	return placeholders("> ?")
}

// =============================================================================
// KeyWords

func (q *queryBuilder) SELECT(fields ...field) *queryBuilder {
	q.b.WriteString("SELECT ")
	q.b.WriteString(q.concatFields(fields...))
	return q
}

func (q *queryBuilder) INSERT_INTO(t table, fields ...field) *queryBuilder {
	q.b.WriteString("INSERT INTO ")
	q.b.WriteString(t.String())
	q.b.WriteString(" (")
	q.b.WriteString(q.concatFields(fields...))
	q.b.WriteString(")")
	return q
}

func (q *queryBuilder) FROM(t table) *queryBuilder {
	q.b.WriteString(" FROM ")
	q.b.WriteString(t.String())
	return q
}

func (q *queryBuilder) JOIN(t table, o on) *queryBuilder {
	q.b.WriteString(" JOIN ")
	q.b.WriteString(t.String())
	q.b.WriteString(" ON ")
	q.b.WriteString(string(o))
	return q
}

func (q *queryBuilder) WHERE(f field, ph placeholders) *queryBuilder {
	q.b.WriteString(" WHERE ")
	q.b.WriteString(f.String())
	q.b.WriteString(" ")
	q.b.WriteString(string(ph))
	return q
}

func (q *queryBuilder) AND(f field, ph placeholders) *queryBuilder {
	q.b.WriteString(" AND ")
	q.b.WriteString(f.String())
	q.b.WriteString(" ")
	q.b.WriteString(string(ph))
	return q
}

func (q *queryBuilder) GROUPBY(f field) *queryBuilder {
	q.b.WriteString(" GROUP BY ")
	q.b.WriteString(f.String())
	return q
}

func (q *queryBuilder) VALUES(v ...any) *queryBuilder {
	lenValues := len(v)
	var b strings.Builder
	b.WriteString(" (")
	b.WriteString(q.generatePlaceholdersString(lenValues))
	b.WriteString(")")
	return q
}

func (q *queryBuilder) ON(fieldA field, fieldB field) on {
	var b strings.Builder
	b.WriteString(fieldA.String())
	b.WriteString("= ")
	b.WriteString(fieldB.String())
	return on(b.String())
}

func (q *queryBuilder) IN(v []any) placeholders {
	lenValues := len(v)
	var b strings.Builder
	b.WriteString("IN (")
	b.WriteString(q.generatePlaceholdersString(lenValues))
	b.WriteString(")")
	return placeholders(b.String())
}

// func (q *queryBuilder) AS(tableName TableName, alias string) QTable {
// 	var b strings.Builder
// 	b.WriteString(string(tableName))
// 	b.WriteString(" AS ")
// 	b.WriteString(alias)
// 	return QTable(b.String())
// }

// func (q *queryBuilder) WITH()

// =============================================================================
// Helpers

func (q *queryBuilder) concatFields(fields ...field) string {
	var b strings.Builder
	for i, f := range fields {
		b.WriteString(f.String())
		if i < len(fields)-1 {
			b.WriteString(", ")
		}
	}
	return b.String()
}

func (q *queryBuilder) generatePlaceholdersString(len int) string {
	p := strings.Repeat("?,", len)
	p = strings.TrimRight(p, ",")
	return p
}
