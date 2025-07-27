package database

import (
	"database/sql"
	"fmt"
	i "server/interfaces"
	"server/services/dbService/SQL/schema"
	prov "server/services/dbService/providers"
	typ "server/types"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(l i.Logger, c prov.Credentials) (*DB, error) {
	conn, err := sql.Open(c.Value("driver"), c.Value("path"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db := &DB{lgr: l, Conn: conn}

	if err := InitDb(db, schema.Get()); err != nil {
		return nil, fmt.Errorf("Database initialization failed: %w", err)
	}

	return db, nil
}

type DB struct {
	lgr  i.Logger
	Conn *sql.DB
}

func (db *DB) Close() {
	db.Conn.Close()
}

func (db *DB) Create(query string, values ...any) (sql.Result, error) {
	res, err := db.Conn.Exec(query, values...)
	if err != nil {
		// var sqliteErr sqlite3.Error
		// if errors.As(err, &sqliteErr) {
		// 	if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		// 		// Handle unique constraint violation here
		// 		fmt.Println("Unique constraint violation!")
		// 	}
		// }
		return nil, fmt.Errorf("db create failed: %w", err)
	}
	return res, nil
}

func (db *DB) Read(query string, values ...any) (typ.Rows, error) {
	rows, err := db.Conn.Query(query, values...)
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
	_, err := db.Conn.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}
	return nil
}

func (db *DB) Delete(query string, conditions ...any) error {
	_, err := db.Conn.Exec(query, conditions...)
	if err != nil {
		return fmt.Errorf("delete failed: %w", err)
	}
	return nil
}
