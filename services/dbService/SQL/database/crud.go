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
		return nil, fmt.Errorf("Failed to connect to database::%w \n args::%v", err, c)
	}

	db := &DB{
		lgr:  l,
		Conn: conn,
	}

	if err := InitDb(db, schema.Get()); err != nil {
		return nil, err
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
	db.lgr.LogFunctionInfo()
	res, err := db.Conn.Exec(query, values...)
	if err != nil {
		// var sqliteErr sqlite3.Error
		// if errors.As(err, &sqliteErr) {
		// 	if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
		// 		return nil, fmt.Errorf(errMessage, errUniqueConstraintFail, err, query, values)
		// 	}
		// }
		return nil, fmt.Errorf("Insert failed::%w \n args::%s; %v", err, query, values)
	}
	return res, nil
}

func (db *DB) Read(query string, values ...any) (typ.Rows, error) {
	db.lgr.LogFunctionInfo()
	rows, err := db.Conn.Query(query, values...)
	if err != nil {
		errRowsNotFound := fmt.Errorf("Read failed::%w \n args::%s; %v", err, query, values)
		db.lgr.LogError(errRowsNotFound)
		return nil, nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("Failed to get columns::%w \n Rows::%v", err, rows)
	}

	var results typ.Rows

	for rows.Next() {
		cols := make([]any, len(columns))
		colPtrs := make([]any, len(columns))
		for i := range cols {
			colPtrs[i] = &cols[i]
		}

		if err := rows.Scan(colPtrs...); err != nil {
			return nil, fmt.Errorf("Row scan failed::%w \n Rows::%v", err, rows)
		}

		rowMap := make(map[string]any)
		for i, col := range columns {
			rowMap[col] = cols[i]
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Row iteration error::%w \n Rows::%v", err, rows)
	}

	return results, nil
}

func (db *DB) Update(query string, values ...any) error {
	db.lgr.LogFunctionInfo()
	_, err := db.Conn.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("Update failed::%w \n args::%s; %v", err, query, values)
	}
	return nil
}

func (db *DB) Delete(query string, values ...any) error {
	db.lgr.LogFunctionInfo()
	_, err := db.Conn.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("Delete failed::%w \n args::%s; %v", err, query, values)
	}
	return nil
}
