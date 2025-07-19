package dbService

import (
	"strings"
)

// The returned QueryBuilder simply appends strings, it is up to the
// user to construct valid SQL queries. Misuse is not guarded against.
func NewQueryBuilder() *queryBuilder {
	qb := &queryBuilder{
		strings.Builder{},
	}
	return qb
}

type queryBuilder struct {
	strings.Builder
}

// type field string
type on string
type placeholders string
type values string

type SQLQuery string

func (s SQLQuery) String() string {
	return string(s)
}

type field string

func (f field) String() string {
	return string(f)
}

type table string

func (t table) String() string {
	return string(t)
}

func (q *queryBuilder) Build() SQLQuery {
	s := strings.TrimSpace(q.String())
	return SQLQuery(s)
}

// =============================================================================
// Objects

func (q *queryBuilder) Field(name string) field {
	return field(name)
}

func (q *queryBuilder) FieldWithAlias(name string, tag string) field {
	var b strings.Builder
	b.WriteString(tag)
	b.WriteString(".")
	b.WriteString(name)
	return field(b.String())
}

func (q *queryBuilder) All() field {
	return field("*")
}

func (q *queryBuilder) Table(name string) table {
	return table(name)
}

func (q *queryBuilder) TableWithAlias(name string, tag string) table {
	var b strings.Builder
	b.WriteString(name)
	b.WriteString(" AS ")
	b.WriteString(tag)
	return table(b.String())
}

// =============================================================================
// Functions

func (q *queryBuilder) Count(f field) field {
	var b strings.Builder
	b.WriteString(" COUNT(")
	b.WriteString(f.String())
	b.WriteString(") AS count")
	return field(b.String())
}

func (q *queryBuilder) Max(f field) field {
	var b strings.Builder
	b.WriteString(" MAX(")
	b.WriteString(f.String())
	b.WriteString(") AS max")
	return field(b.String())
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
	q.WriteString("SELECT ")
	q.WriteString(q.concatFields(fields...))
	return q
}

func (q *queryBuilder) INSERT_INTO(t table, fields ...field) *queryBuilder {
	q.WriteString("INSERT INTO ")
	q.WriteString(t.String())
	q.WriteString(" (")
	q.WriteString(q.concatFields(fields...))
	q.WriteString(")")
	return q
}

func (q *queryBuilder) FROM(t table) *queryBuilder {
	q.WriteString(" FROM ")
	q.WriteString(t.String())
	return q
}

func (q *queryBuilder) JOIN(t table, o on) *queryBuilder {
	q.WriteString(" JOIN ")
	q.WriteString(t.String())
	q.WriteString(" ON ")
	q.WriteString(string(o))
	return q
}

func (q *queryBuilder) WHERE(f field, ph placeholders) *queryBuilder {
	q.WriteString(" WHERE ")
	q.WriteString(f.String())
	q.WriteString(" ")
	q.WriteString(string(ph))
	return q
}

func (q *queryBuilder) AND(f field, ph placeholders) *queryBuilder {
	q.WriteString(" AND ")
	q.WriteString(f.String())
	q.WriteString(" ")
	q.WriteString(string(ph))
	return q
}

func (q *queryBuilder) GROUPBY(f field) *queryBuilder {
	q.WriteString(" GROUP BY ")
	q.WriteString(f.String())
	return q
}

func (q *queryBuilder) VALUES(v ...any) *queryBuilder {
	lenValues := len(v)
	q.WriteString(" VALUES ")
	q.WriteString("(")
	q.WriteString(q.generatePlaceholdersString(lenValues))
	q.WriteString(")")
	return q
}

func (q *queryBuilder) ON(fieldA field, fieldB field) on {
	var b strings.Builder
	b.WriteString(fieldA.String())
	b.WriteString(" = ")
	b.WriteString(fieldB.String())
	return on(b.String())
}

func (q *queryBuilder) IN(v ...any) placeholders {
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
	p := strings.Repeat("?, ", len)
	p = strings.TrimRight(p, ", ")
	return p
}
