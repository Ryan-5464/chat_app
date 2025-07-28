package database

import "fmt"

const (
	errDatabaseInitFail string = "database initialization failed:"
)

func InitDb(db *DB, schema []string) error {

	for _, query := range schema {
		_, err := db.Create(query)
		if err != nil {
			return fmt.Errorf("%s %v", errDatabaseInitFail, err)
		}
	}

	return nil
}
