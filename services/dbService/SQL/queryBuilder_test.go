package SQL

import (
	"testing"
)

func TestField(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.Field("test")
	expected := "test"

	if result.String() != expected {
		t.Errorf("Field() = %s; want %s", result, expected)
	}
}

func TestFieldWithAlias(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.FieldWithAlias("test", "t")
	expected := "t.test"

	if result.String() != expected {
		t.Errorf("FieldWithAlias() = %s; want %s", result, expected)
	}
}

func TestAll(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.All()
	expected := "*"

	if result.String() != expected {
		t.Errorf("All() = %s; want %s", result, expected)
	}
}

func TestTable(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.Table("testTable")
	expected := "testTable"

	if result.String() != expected {
		t.Errorf("Table() = %s; want %s", result, expected)
	}
}

func TestTableWithAlias(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.TableWithAlias("testTable", "t")
	expected := "testTable AS t"

	if result.String() != expected {
		t.Errorf("TableWithAlias() = %s; want %s", result, expected)
	}
}

func TestCount(t *testing.T) {
	qb := NewQueryBuilder()

	f := qb.Field("test")

	result := qb.Count(f)
	expected := " COUNT(test) AS count"

	if result.String() != expected {
		t.Errorf("Count() = %s; want %s", result, expected)
	}

}

func TestMax(t *testing.T) {
	qb := NewQueryBuilder()

	f := qb.Field("test")

	result := qb.Max(f)
	expected := " MAX(test) AS max"

	if result.String() != expected {
		t.Errorf("Max() = %s; want %s", result, expected)
	}
}

func TestGreaterThanPlaceholder(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.GreaterThan()
	expected := "> ?"

	if string(result) != expected {
		t.Errorf("GreaterThan() = %s; want %s", result, expected)
	}
}

func TestINWithOneValue(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.IN(1)
	expected := "IN (?)"

	if string(result) != expected {
		t.Errorf("IN() with 1 value = %s; want %s", result, expected)
	}
}

func TestINWithMultipleValues(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.IN(1, 2)
	expected := "IN (?, ?)"

	if string(result) != expected {
		t.Errorf("IN() with multiple values = %s; want %s", result, expected)
	}
}

func TestSimpleSELECTQuery(t *testing.T) {
	qb := NewQueryBuilder()

	f1 := qb.Field("f1")
	f2 := qb.Field("f2")
	f3 := qb.Field("f3")
	tab := qb.Table("test")

	result := qb.SELECT(f1, f2, f3).
		FROM(tab).
		WHERE(f1, qb.EqualTo()).
		Build()
	expected := "SELECT f1, f2, f3 FROM test WHERE f1 = ?"

	if result.String() != expected {
		t.Errorf("result = %s; want %s", result, expected)
	}
}

func TestINSERT_INTOQuery(t *testing.T) {
	qb := NewQueryBuilder()

	f1 := qb.Field("f1")
	f2 := qb.Field("f2")
	f3 := qb.Field("f3")
	tab := qb.Table("test")
	v1 := 1
	v2 := 2
	v3 := 3

	result := qb.INSERT_INTO(tab, f1, f2, f3).
		VALUES(v1, v2, v3).
		Build()

	expected := "INSERT INTO test (f1, f2, f3) VALUES (?, ?, ?)"

	if result.String() != expected {
		t.Errorf("result = %s; want %s", result, expected)
	}
}

func TestComplexJOINQuery(t *testing.T) {
	qb := NewQueryBuilder()

	tag1 := "a"
	tag2 := "b"
	f1 := qb.FieldWithAlias("f1", tag1)
	f2 := qb.FieldWithAlias("f2", tag1)
	f3 := qb.FieldWithAlias("f3", tag2)
	tab1 := qb.TableWithAlias("test1", tag1)
	tab2 := qb.TableWithAlias("test2", tag2)
	v1 := 1
	v2 := 2

	result := qb.SELECT(qb.All()).
		FROM(tab1).
		JOIN(tab2, qb.ON(f1, f3)).
		WHERE(f2, qb.IN(v1, v2)).
		AND(f2, qb.EqualTo()).
		GROUPBY(f2).
		Build()

	expected := `SELECT * FROM test1 AS a JOIN test2 AS b ON a.f1 = b.f3 WHERE a.f2 IN (?, ?) AND a.f2 = ? GROUP BY a.f2`

	if result.String() != expected {
		t.Errorf("result = %s; want \n %s", result, expected)
	}

}
