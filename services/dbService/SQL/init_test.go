package SQL

import (
	"database/sql"
	"fmt"
	schema "server/services/dbService/SQL/schema"
	prov "server/services/dbService/providers"
	"testing"
)

func TestDBInitialization(t *testing.T) {
	config := prov.DbConfig{}
	config.Add("driver", "sqlite3")
	config.Add("path", ":memory:")
	c := prov.NewCredentials(prov.SQLite3, config)
	dbs, err := NewDbService(nil, c)
	if err != nil {
		t.Fatalf("failed to initialize database: %v", err)
	}
	defer dbs.Close()

	expectedTables := []string{
		schema.UserTable,
		schema.MessageTable,
		schema.ChatTable,
		schema.MemberTable,
	}

	ok, err := tablesExist(dbs.db.Conn, expectedTables)
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
