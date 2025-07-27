package database

func InitDb(db *DB, schema []string) error {

	for _, query := range schema {
		_, err := db.Create(query)
		if err != nil {
			return err
		}
	}

	return nil
}
