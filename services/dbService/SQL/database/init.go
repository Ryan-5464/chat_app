package database

import (
	"fmt"
	"log"
)

const (
	errDatabaseInitFail string = "database initialization failed:"
)

func InitDb(db *DB, schema []string) error {

	for _, query := range schema {
		log.Println(query)
		_, err := db.Create(query)
		if err != nil {
			return fmt.Errorf("%s %v", errDatabaseInitFail, err)
		}
	}

	return nil
}
