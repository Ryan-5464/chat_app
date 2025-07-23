package database

import (
	"database/sql"
	"fmt"
	typ "server/types"
)

func NewDb(conn *sql.DB) *DB {
	return &DB{conn: conn}
}

type DB struct {
	conn *sql.DB
}

func (db *DB) Create(query string, values ...any) (sql.Result, error) {
	res, err := db.conn.Exec(query, values...)
	if err != nil {
		return nil, fmt.Errorf("db create failed: %w", err)
	}
	return res, nil
}

func (db *DB) Read(query string, values ...any) (typ.Rows, error) {
	rows, err := db.conn.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("db read failed: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results typ.Rows

	for rows.Next() {
		cols := make([]any, len(columns))
		colPtrs := make([]any, len(columns))
		for i := range cols {
			colPtrs[i] = &cols[i]
		}

		if err := rows.Scan(colPtrs...); err != nil {
			return nil, fmt.Errorf("row scan failed: %w", err)
		}

		rowMap := make(map[string]any)
		for i, col := range columns {
			rowMap[col] = cols[i]
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return results, nil
}

func (db *DB) Update(query string, values ...any) error {
	_, err := db.conn.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

func (db *DB) Delete(query string, conditions ...any) error {
	_, err := db.conn.Exec(query, conditions...)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}
