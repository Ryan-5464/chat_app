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

func (dbs *DbService) GetUser(usrId typ.UserId) (model.User, error) {
	usrIds := []typ.UserId{usrId}

	usrs, err := dbs.GetUsers(usrIds)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user from database: %w", err)
	}

	return usrs[0], nil
}

func (dbs *DbService) GetUsers(usrIds []typ.UserId) ([]model.User, error) {
	qb := qbuilder.NewQueryBuilder()

	usrTbl := qb.Table(schema.UserTable)
	usrIdF := qb.Field(schema.UserId)

	ids := ToAnySlice(usrIds)

	query := qb.SELECT(qb.All()).FROM(usrTbl).WHERE(usrIdF, qb.IN(ids))

	rows, err := dbs.db.Read(query.String(), ids)
	if err != nil {
		return nil, fmt.Errorf("failed to read users from database: %w", err)
	}

	usrMs := populateUserModels(rows)

	return usrMs, err
}

func (dbs *DbService) GetChats() []entities.Chat {
	return []entities.Chat{}
}

func (dbs *DbService) GetMessage(msgId typ.MessageId) (model.Message, error) {
	msgIds := []typ.MessageId{msgId}

	msgs, err := dbs.GetMessages(msgIds)
	if err != nil {
		return model.Message{}, err
	}

	return msgs[0], nil
}

func (dbs *DbService) GetMessages(msgIds []typ.MessageId) ([]model.Message, error) {
	qb := qbuilder.NewQueryBuilder()

	msgTbl := qb.Table(schema.MessageTable)
	msgIdF := qb.Field(schema.MessageId)

	ids := ToAnySlice(msgIds)

	query := qb.SELECT(qb.All()).FROM(msgTbl).WHERE(msgIdF, qb.IN(ids...))

	rows, err := dbs.db.Read(query.String(), ids...)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages from database: %w", err)
	}

	return populateMessageModels(rows), nil
}

func (dbs *DbService) NewUser(u model.User) (model.User, error) {
	qb := qbuilder.NewQueryBuilder()

	usrTbl := qb.Table(schema.UserTable)
	usrIdF := qb.Field(schema.UserId)
	usrNameF := qb.Field(schema.Name)
	emailF := qb.Field(schema.Email)
	pwdHashF := qb.Field(schema.PwdHash)
	joinedF := qb.Field(schema.Joined)

	query := qb.INSERT_INTO(
		usrTbl,
		usrIdF,
		usrNameF,
		emailF,
		pwdHashF,
		joinedF,
	).VALUES(
		u.Id,
		u.Name,
		u.Email,
		u.PwdHash,
		u.Joined,
	)

	res, err := dbs.db.Create(
		query.String(),
		u.Id,
		u.Name,
		u.Email,
		u.PwdHash,
		u.Joined,
	)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to create user in database: %w", err)
	}

	newUsrId, err := res.LastInsertId()
	if err != nil {
		return model.User{}, err
	}

	usr, err := dbs.GetUser(typ.UserId(newUsrId))
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get message from database: %w", err)
	}

	return usr, nil
}

func (dbs *DbService) NewMessage(m model.Message) (model.Message, error) {
	qb := qbuilder.NewQueryBuilder()

	msgTbl := qb.Table(schema.MessageTable)
	msgIdF := qb.Field(schema.MessageId)
	usrIdF := qb.Field(schema.UserId)
	chatIdF := qb.Field(schema.ChatId)
	replyIdF := qb.Field(schema.ReplyId)
	msgTextF := qb.Field(schema.MsgText)
	createdAtF := qb.Field(schema.CreatedAt)
	lastEditAtF := qb.Field(schema.LastEditAt)

	query := qb.INSERT_INTO(
		msgTbl,
		msgIdF,
		usrIdF,
		chatIdF,
		replyIdF,
		msgTextF,
		createdAtF,
		lastEditAtF,
	).VALUES(
		m.Id,
		m.UserId,
		m.ChatId,
		m.ReplyId,
		m.Text,
		m.CreatedAt,
		m.LastEditAt,
	)

	res, err := dbs.db.Create(
		query.String(),
		m.Id,
		m.UserId,
		m.ChatId,
		m.ReplyId,
		m.Text,
		m.CreatedAt,
		m.LastEditAt,
	)
	if err != nil {
		return model.Message{}, fmt.Errorf("failed to create message in database: %w", err)
	}

	newMsgId, err := res.LastInsertId()
	if err != nil {
		return model.Message{}, err
	}

	msg, err := dbs.GetMessage(typ.MessageId(newMsgId))
	if err != nil {
		return model.Message{}, fmt.Errorf("failed to get message from database: %w", err)
	}

	return msg, nil
}

func populateMessageModels(rows typ.Rows) []model.Message {

	var msgMs []model.Message
	for _, row := range rows {
		msgM := model.Message{
			Id:         typ.MessageId(row[schema.MessageId].(int64)),
			UserId:     typ.UserId(row[schema.UserId].(int64)),
			ChatId:     typ.ChatId(row[schema.ChatId].(int64)),
			ReplyId:    typ.MessageId(row[schema.ReplyId].(int64)),
			Text:       row[schema.MsgText].(string),
			CreatedAt:  row[schema.CreatedAt].(time.Time),
			LastEditAt: row[schema.LastEditAt].(time.Time),
		}
		msgMs = append(msgMs, msgM)
	}

	return msgMs
}

func populateUserModels(rows typ.Rows) []model.User {
	var usrMs []model.User
	for _, row := range rows {
		usrM := model.User{
			Id:      typ.UserId(row[schema.UserId].(int64)),
			Name:    row[schema.Name].(string),
			Email:   cred.Email(row[schema.Email].(string)),
			PwdHash: cred.PwdHash(row[schema.PwdHash].(string)),
			Joined:  row[schema.CreatedAt].(time.Time),
		}
		usrMs = append(usrMs, usrM)
	}

	return usrMs
}

func ConvertSlice[T any](input []any) ([]T, error) {
	result := make([]T, len(input))
	for i, v := range input {
		val, ok := v.(T)
		if !ok {
			return nil, fmt.Errorf("element at index %d is not of type %T", i, val)
		}
		result[i] = val
	}
	return result, nil
}

func ToAnySlice[T any](input []T) []any {
	result := make([]any, len(input))
	for i, v := range input {
		result[i] = v
	}
	return result
}
