package SQL

import (
	"database/sql"
	"fmt"
	"server/data/entities"
	cred "server/services/authService/credentials"
	d "server/services/dbService/SQL/database"
	model "server/services/dbService/SQL/models"
	qbuilder "server/services/dbService/SQL/querybuilder"
	schema "server/services/dbService/SQL/schema"
	prov "server/services/dbService/providers"
	typ "server/types"
	"time"

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

func (dbs *DbService) Close() {
	dbs.db.Close()
}

func (dbs *DbService) GetUsers(chatId typ.ChatId) ([]model.User, error) {
	qb := qbuilder.NewQueryBuilder()

	usrTbl := qb.Table(schema.UserTable)
	chatIdF := qb.Field(schema.ChatId)

	query := qb.SELECT(qb.All()).FROM(usrTbl).WHERE(chatIdF, qb.EqualTo())

	rows, err := dbs.db.Read(query.String(), chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to read users from database: %w", err)
	}

	var mdls []model.User
	for _, row := range rows {
		u := model.User{}
		u.Id = typ.UserId(row[schema.UserId].(int64))
		u.Name = row[schema.Name].(string)
		u.Email = cred.Email(row[schema.Email].(string))
		u.PwdHash = cred.PwdHash(row[schema.PwdHash].(string))
		u.Joined = row[schema.CreatedAt].(time.Time)
		mdls = append(mdls, u)
	}

	return []model.User{}, err
}

func (dbs *DbService) GetChats() []entities.Chat {
	return []entities.Chat{}
}

func (dbs *DbService) GetMessages() []entities.Message {
	return []entities.Message{}
}
