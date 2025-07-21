package SQL

import (
	"database/sql"
	"fmt"
	"server/data/entities"
	prov "server/services/dbService/providers"

	_ "github.com/mattn/go-sqlite3"
)

func NewDbService(c prov.Credentials) (*DbService, error) {
	conn, err := sql.Open(c.Value("driver"), c.Value("path"))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	dbS := &DbService{db: newDb(conn)}
	return dbS, nil
}

type DbService struct {
	db *database
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
