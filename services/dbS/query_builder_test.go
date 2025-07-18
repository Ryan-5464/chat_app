package dbService

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
		t.Errorf("Field() = %s; want %s", result, expected)
	}
}

func TestAll(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.All()
	expected := " * "

	if result.String() != expected {
		t.Errorf("Field() = %s; want %s", result, expected)
	}
}

func TestTable(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.Table("testTable")
	expected := "testTable"

	if result.String() != expected {
		t.Errorf("Field() = %s; want %s", result, expected)
	}
}

func TestTableWithAlias(t *testing.T) {
	qb := NewQueryBuilder()
	result := qb.TableWithAlias("testTable", "t")
	expected := "testTable AS t"

	if result.String() != expected {
		t.Errorf("Field() = %s; want %s", result, expected)
	}
}

func TestCount(t *testing.T) {
	qb := NewQueryBuilder()

	f := qb.Field("test")

	result := qb.Count(f)
	expected := " COUNT(test) AS count"

	if result.String() != expected {
		t.Errorf("Field() = %s; want %s", result, expected)
	}

}
