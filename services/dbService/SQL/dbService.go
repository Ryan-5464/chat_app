package SQL

import (
	"database/sql"
	"fmt"
	"server/data/entities"
	d "server/services/dbService/SQL/database"
	"server/services/dbService/SQL/schema"
	prov "server/services/dbService/providers"

	_ "github.com/mattn/go-sqlite3"
)

func NewDbService(c prov.Credentials) (*DbService, error) {
	conn, err := sql.Open(c.Value("driver"), c.Value("path"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	db := d.NewDb(conn)
	initDb(db, schema.Get())
	dbS := &DbService{db: db}
	return dbS, nil
}

type DbService struct {
	db *d.DB
}

func (db *DbService) GetUsers() []entities.User {
	return []entities.User{}
}

func (db *DbService) GetChats() []entities.Chat {
	return []entities.Chat{}
}

func (db *DbService) GetMessages() []entities.Message {
	return []entities.Message{}
}
