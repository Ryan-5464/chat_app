package SQL

import (
	d "server/services/dbService/SQL/database"
)

func initDb(db *d.DB, schema []string) error {

	for _, query := range schema {
		_, err := db.Create(query)
		if err != nil {
			return err
		}
	}

	return nil
}
