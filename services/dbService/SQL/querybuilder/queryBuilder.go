package qbuilder

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

type Field string

func (f Field) String() string {
	return string(f)
}

type table string

func (t table) String() string {
	return string(t)
}

func (q *queryBuilder) Build() string {
	s := strings.TrimSpace(q.String())
	return s
}

// =============================================================================
// Objects

func (q *queryBuilder) Field(name string) Field {
	return Field(name)
}

func (q *queryBuilder) FieldWithAlias(name string, tag string) Field {
	var b strings.Builder
	b.WriteString(tag)
	b.WriteString(".")
	b.WriteString(name)
	return Field(b.String())
}

func (q *queryBuilder) All() Field {
	return Field("*")
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

func (q *queryBuilder) Count(f Field) Field {
	var b strings.Builder
	b.WriteString(" COUNT(")
	b.WriteString(f.String())
	b.WriteString(") AS count")
	return Field(b.String())
}

func (q *queryBuilder) Max(f Field) Field {
	var b strings.Builder
	b.WriteString(" MAX(")
	b.WriteString(f.String())
	b.WriteString(") AS max")
	return Field(b.String())
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

func (q *queryBuilder) SELECT(fields ...Field) *queryBuilder {
	q.WriteString("SELECT ")
	q.WriteString(q.concatFields(fields...))
	return q
}

func (q *queryBuilder) DELETE_FROM(t table) *queryBuilder {
	q.WriteString("DELETE FROM ")
	q.WriteString(t.String())
	return q
}

func (q *queryBuilder) INSERT_INTO(t table, fields ...Field) *queryBuilder {
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

func (q *queryBuilder) WHERE(f Field, ph placeholders) *queryBuilder {
	q.WriteString(" WHERE ")
	q.WriteString(f.String())
	q.WriteString(" ")
	q.WriteString(string(ph))
	return q
}

func (q *queryBuilder) OR(f Field, ph placeholders) *queryBuilder {
	q.WriteString(" OR ")
	q.WriteString(f.String())
	q.WriteString(" ")
	q.WriteString(string(ph))
	return q
}

func (q *queryBuilder) AND(f Field, ph placeholders) *queryBuilder {
	q.WriteString(" AND ")
	q.WriteString(f.String())
	q.WriteString(" ")
	q.WriteString(string(ph))
	return q
}

func (q *queryBuilder) GROUPBY(f Field) *queryBuilder {
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

func (q *queryBuilder) ON(fieldA Field, fieldB Field) on {
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

// =============================================================================
// Helpers

func (q *queryBuilder) concatFields(fields ...Field) string {
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
