package database

import (
	"database/sql"
	"fmt"
	"server/logging"
	"server/services/dbService/SQL/schema"
	"server/services/dbService/providers"
	"testing"
)

func TestDBInitialization_Success(t *testing.T) {

	config := providers.DbConfig{}
	config.Add("driver", providers.SQLite3.String())
	config.Add("path", providers.InMemoryDb.String())
	dbCredentials := providers.NewDbCredentials(providers.SQLite3, config)
	logger := logging.NewLogger(false)

	db, err := NewDatabase(logger, dbCredentials)
	if err != nil {
		t.Fatalf("Database creation failed: %v", err)
	}
	defer db.Close()

	expectedTables := []string{
		schema.UserTable,
		schema.MessageTable,
		schema.ChatTable,
		schema.MemberTable,
	}

	ok, err := tablesExist(db.Conn, expectedTables)
	if err != nil {
		t.Fatalf("Database initialization failed: %v", err)
	}
	if !ok {
		t.Fatal("Not all expected tables exist")
	}
}

func tablesExist(db *sql.DB, tableNames []string) (bool, error) {
	for _, table := range tableNames {
		var name string
		err := db.QueryRow(`
			SELECT name FROM sqlite_master 
			WHERE type='table' AND name=?;
		`, table).Scan(&name)

		if err == sql.ErrNoRows {
			return false, fmt.Errorf("table %s does not exist", table)
		} else if err != nil {
			return false, err
		}
	}
	return true, nil
}
