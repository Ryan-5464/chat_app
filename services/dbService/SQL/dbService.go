package SQL

import (
	"database/sql"
	"fmt"
	"server/data/entities"
	typ "server/types"

	_ "github.com/mattn/go-sqlite3"
)

func NewDbService(c typ.Credentials) (*DbService, error) {
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
