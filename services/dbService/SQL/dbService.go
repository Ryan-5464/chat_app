package SQL

import (
	"errors"
	"fmt"
	"log"
	i "server/interfaces"
	cred "server/services/authService/credentials"
	"server/services/dbService/SQL/database"
	model "server/services/dbService/SQL/models"
	qbuilder "server/services/dbService/SQL/querybuilder"
	schema "server/services/dbService/SQL/schema"
	prov "server/services/dbService/providers"
	typ "server/types"
	"time"
)

func NewDbService(lgr i.Logger, c prov.Credentials) (*DbService, error) {
	db, err := database.NewDatabase(lgr, c)
	if err != nil {
		return nil, err
	}

	dbS := &DbService{
		lgr: lgr,
		db:  db,
	}
	return dbS, nil
}

type DbService struct {
	lgr i.Logger
	db  *database.DB
}

func (dbs *DbService) Close() {
	dbs.db.Close()
}

func (dbs *DbService) GetChat(chatId typ.ChatId) (model.Chat, error) {
	dbs.lgr.LogFunctionInfo()

	qb := qbuilder.NewQueryBuilder()

	chatTbl := qb.Table(schema.ChatTable)
	chatIdF := qb.Field(schema.ChatId)

	query := qb.SELECT(qb.All()).FROM(chatTbl).WHERE(chatIdF, qb.EqualTo())

	rows, err := dbs.db.Read(query.String(), chatId)
	if err != nil {
		return model.Chat{}, err
	}

	if len(rows) == 0 {
		return model.Chat{}, nil
	}

	chatMs := populateChatModels(rows)

	return chatMs[0], nil

}

func (dbs *DbService) GetChats(userId typ.UserId) ([]model.Chat, error) {
	dbs.lgr.LogFunctionInfo()

	qb := qbuilder.NewQueryBuilder()

	chatTbl := qb.Table(schema.ChatTable)
	adminIdF := qb.Field(schema.AdminId)

	query := qb.SELECT(qb.All()).FROM(chatTbl).WHERE(adminIdF, qb.EqualTo())

	rows, err := dbs.db.Read(query.String(), userId)
	if err != nil {
		return []model.Chat{}, fmt.Errorf("failed to read chats from database")
	}
	log.Println("ROWS : ", rows)
	chatMs := populateChatModels(rows)

	return chatMs, nil
}

func (dbs *DbService) FindUser(email cred.Email) (model.User, error) {
	dbs.lgr.LogFunctionInfo()

	qb := qbuilder.NewQueryBuilder()

	usrTbl := qb.Table(schema.UserTable)
	emailF := qb.Field(schema.Email)

	query := qb.SELECT(qb.All()).FROM(usrTbl).WHERE(emailF, qb.EqualTo())

	rows, err := dbs.db.Read(query.String(), email)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to find user email from database: %w", err)
	}

	if len(rows) == 0 {
		return model.User{}, nil
	}

	usrs := populateUserModels(rows)

	return usrs[0], err
}

func (dbs *DbService) GetUser(usrId typ.UserId) (model.User, error) {
	dbs.lgr.LogFunctionInfo()
	usrIds := []typ.UserId{usrId}

	usrs, err := dbs.GetUsers(usrIds)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user from database: %w", err)
	}

	return usrs[0], nil
}

func (dbs *DbService) GetUsers(usrIds []typ.UserId) ([]model.User, error) {
	dbs.lgr.LogFunctionInfo()
	qb := qbuilder.NewQueryBuilder()

	usrTbl := qb.Table(schema.UserTable)
	usrIdF := qb.Field(schema.UserId)

	ids := ToAnySlice(usrIds)

	query := qb.SELECT(qb.All()).FROM(usrTbl).WHERE(usrIdF, qb.IN(ids))

	rows, err := dbs.db.Read(query.String(), ids...)
	if err != nil {
		return nil, fmt.Errorf("failed to read users from database: %w", err)
	}

	usrMs := populateUserModels(rows)

	return usrMs, err
}

func (dbs *DbService) GetMembers(chatId typ.ChatId) ([]model.Member, error) {
	dbs.lgr.LogFunctionInfo()
	qb := qbuilder.NewQueryBuilder()

	mbrTbl := qb.Table(schema.MemberTable)
	chatIdF := qb.Field(schema.ChatId)

	query := qb.SELECT(qb.All()).FROM(mbrTbl).WHERE(chatIdF, qb.EqualTo())
	rows, err := dbs.db.Read(query.String(), chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to read members from database: %w", err)
	}

	mbrMs := populateMemberModels(rows)

	return mbrMs, err
}

func (dbs *DbService) GetUsersForChat(chatId typ.ChatId) ([]model.User, error) {
	dbs.lgr.LogFunctionInfo()

	mbrMs, err := dbs.GetMembers(chatId)
	if err != nil {
		return nil, err
	}

	var usrIds []typ.UserId
	for _, mdl := range mbrMs {
		usrIds = append(usrIds, mdl.UserId)
	}

	usrMs, err := dbs.GetUsers(usrIds)
	if err != nil {
		return nil, err
	}

	return usrMs, nil

}

func (dbs *DbService) GetMessage(msgId typ.MessageId) (model.Message, error) {
	dbs.lgr.LogFunctionInfo()
	msgIds := []typ.MessageId{msgId}

	msgs, err := dbs.GetMessages(msgIds)
	if err != nil {
		return model.Message{}, err
	}

	return msgs[0], nil
}

func (dbs *DbService) GetChatMessages(chatId typ.ChatId) ([]model.Message, error) {
	dbs.lgr.LogFunctionInfo()
	qb := qbuilder.NewQueryBuilder()

	msgTbl := qb.Table(schema.MessageTable)
	chatIdF := qb.Field(schema.ChatId)

	query := qb.SELECT(qb.All()).FROM(msgTbl).WHERE(chatIdF, qb.EqualTo())

	rows, err := dbs.db.Read(query.String(), chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to read messages from database: %w", err)
	}

	return populateMessageModels(rows), nil
}

func (dbs *DbService) GetMessages(msgIds []typ.MessageId) ([]model.Message, error) {
	dbs.lgr.LogFunctionInfo()
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

func (dbs *DbService) NewChat(c model.Chat) (model.Chat, error) {
	dbs.lgr.LogFunctionInfo()
	qb := qbuilder.NewQueryBuilder()

	chatTbl := qb.Table(schema.ChatTable)
	chatName := qb.Field(schema.Name)
	adminId := qb.Field(schema.AdminId)

	query := qb.INSERT_INTO(chatTbl, chatName, adminId).
		VALUES(chatName, adminId)

	res, err := dbs.db.Create(query.String(), c.Name, c.AdminId)
	if err != nil {
		return model.Chat{}, fmt.Errorf("failed to create chat in database: %w", err)
	}

	newChatId, err := res.LastInsertId()
	if err != nil {
		return model.Chat{}, err
	}

	chat, err := dbs.GetChat(typ.ChatId(newChatId))
	if err != nil {
		return model.Chat{}, fmt.Errorf("faield to get chat")
	}

	return chat, nil
}

func (dbs *DbService) NewMember(chatId typ.ChatId, userId typ.UserId) error {
	dbs.lgr.LogFunctionInfo()
	qb := qbuilder.NewQueryBuilder()

	mbrTbl := qb.Table(schema.MemberTable)
	userIdF := qb.Field(schema.UserId)
	chatIdF := qb.Field(schema.ChatId)

	query := qb.INSERT_INTO(mbrTbl, chatIdF, userIdF).
		VALUES(chatId, userId)
	log.Println(query)

	_, err := dbs.db.Create(query.String(), chatId, userId)
	if err != nil {
		return fmt.Errorf("failed to create member in database: %w", err)
	}

	return nil
}

func (dbs *DbService) NewUser(u model.User) (model.User, error) {
	dbs.lgr.LogFunctionInfo()
	qb := qbuilder.NewQueryBuilder()

	usrTbl := qb.Table(schema.UserTable)
	usrNameF := qb.Field(schema.Name)
	emailF := qb.Field(schema.Email)
	pwdHashF := qb.Field(schema.PwdHash)

	query := qb.INSERT_INTO(usrTbl, usrNameF, emailF, pwdHashF).
		VALUES(u.Name, u.Email, u.PwdHash)
	log.Println(query)
	log.Println(u.Name, u.Email, u.PwdHash)

	res, err := dbs.db.Create(query.String(), u.Name, u.Email, u.PwdHash)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to create user in database: %w", err)
	}

	newUsrId, err := res.LastInsertId()
	if err != nil {
		return model.User{}, err
	}

	usr, err := dbs.GetUser(typ.UserId(newUsrId))
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user from database: %w", err)
	}

	return usr, nil
}

func (dbs *DbService) NewMessage(m model.Message) (model.Message, error) {
	dbs.lgr.LogFunctionInfo()
	qb := qbuilder.NewQueryBuilder()

	msgTbl := qb.Table(schema.MessageTable)
	usrIdF := qb.Field(schema.UserId)
	chatIdF := qb.Field(schema.ChatId)
	replyIdF := qb.Field(schema.ReplyId)
	msgTextF := qb.Field(schema.MsgText)

	query := qb.INSERT_INTO(
		msgTbl,
		usrIdF,
		chatIdF,
		replyIdF,
		msgTextF,
	).VALUES(
		m.UserId,
		m.ChatId,
		m.ReplyId,
		m.Text,
	)

	res, err := dbs.db.Create(
		query.String(),
		m.UserId,
		m.ChatId,
		m.ReplyId,
		m.Text,
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

func (dbs *DbService) FindUsers(e []cred.Email) ([]model.User, error) {
	dbs.lgr.LogFunctionInfo()

	var users []model.User

	qb := qbuilder.NewQueryBuilder()

	userTable := qb.Table(schema.UserTable)
	emailF := qb.Field(schema.Email)

	emails := ToAnySlice(e)

	query := qb.SELECT(qb.All()).FROM(userTable).WHERE(emailF, qb.IN(emails...))

	rows, err := dbs.db.Read(query.String(), emails...)
	if err != nil {
		return users, err
	}

	if len(rows) == 0 {
		return users, nil
	}

	users = populateUserModels(rows)
	return users, nil
}

func (dbs *DbService) AddContactRelation(userId typ.UserId, contactId typ.UserId) ([]model.ContactRelation, error) {
	dbs.lgr.LogFunctionInfo()

	var contactRelations []model.ContactRelation

	qb := qbuilder.NewQueryBuilder()

	contactsTable := qb.Table(schema.ContactsTable)
	contact1F := qb.Field(schema.Contact1)
	contact2F := qb.Field(schema.Contact2)

	query := qb.INSERT_INTO(contactsTable, contact1F, contact2F).
		VALUES(c.Contact1, c.Contact2)

	res, err := dbs.db.Create(query.String(), userId, contactId)
	if err != nil {
		return contactRelations, err
	}

	return dbs.getContactRelation(userId, res.LastInsertId())
}

// Specific helper function for AddContactRelation.
func (dbs *DbService) getContactRelation(userId typ.UserId, rowId int64, err error) ([]model.ContactRelation, error) {
	dbs.lgr.LogFunctionInfo()

	var contactRelations []model.ContactRelation

	qb := qbuilder.NewQueryBuilder()

	contactsTable := qb.Table(schema.ContactsTable)
	rowIdF := qb.Field(schema.RowId)

	query := qb.SELECT(qb.All()).FROM(contactsTable).
		WHERE(rowIdF, qb.EqualTo())

	rows, err := dbs.db.Read(query.String(), rowId)
	if err != nil {
		return contactRelations, err
	}

	if len(rows) == 0 {
		return contactRelations, errors.New("failed to retrieve newly created contact relation")
	}

	return populateContactRelationModels(rows, userId), nil
}

func (dbs *DbService) GetContactRelations(userId typ.UserId) ([]model.ContactRelation, error) {
	dbs.lgr.LogFunctionInfo()

	var contactRelations []model.ContactRelation

	qb := qbuilder.NewQueryBuilder()

	contactsTable := qb.Table(schema.ContactsTable)
	contact1F := qb.Field(schema.Contact1)
	contact2F := qb.Field(schema.Contact2)

	query := qb.SELECT(qb.All()).FROM(contactsTable).
		WHERE(contact1F, qb.EqualTo()).OR(contact2F, qb.EqualTo())

	rows, err := dbs.db.Read(query.String(), userId, userId)
	if err != nil {
		return contactRelations, err
	}

	if len(rows) == 0 {
		return contactRelations, nil
	}

	return populateContactRelationModels(rows, userId), nil
}

func populateChatModels(rows typ.Rows) []model.Chat {
	log.Println("test")
	var chatMs []model.Chat
	for _, row := range rows {
		chatM := model.Chat{
			Id:        parseChatId(row[schema.ChatId]),
			Name:      parseString(row[schema.Name]),
			AdminId:   parseUserId(row[schema.AdminId]),
			CreatedAt: parseTime(row[schema.CreatedAt]),
		}
		chatMs = append(chatMs, chatM)
	}
	return chatMs
}

func populateMessageModels(rows typ.Rows) []model.Message {

	var msgMs []model.Message
	for _, row := range rows {
		msgM := model.Message{
			Id:         parseMessageId(row[schema.MessageId]),
			UserId:     parseUserId(row[schema.UserId]),
			ChatId:     parseChatId(row[schema.ChatId]),
			ReplyId:    parseMessageId(row[schema.ReplyId]),
			Text:       parseString(row[schema.MsgText]),
			CreatedAt:  parseTime(row[schema.CreatedAt]),
			LastEditAt: parseTime(row[schema.LastEditAt]),
		}
		msgMs = append(msgMs, msgM)
	}

	return msgMs
}

func populateContactRelationModels(rows typ.Rows, userId typ.UserId) []model.ContactRelation {
	var relations []model.ContactRelation
	for _, row := range rows {
		contactA := parseUserId(row[schema.Contact1])
		contactB := parseUserId(row[schema.Contact2])

		// The user may be in any of the columns in the contactRelation table, so we want to
		// make sure the user is always contact1 in the models so we can assume it in the rest
		// of the program.
		var contact1 typ.UserId
		var contact2 typ.UserId
		if contactA == userId {
			contact1 = contactA
			contact2 = contactB
		} else {
			contact1 = contactB
			contact2 = contactA
		}

		relation := model.ContactRelation{
			Contact1:    contact1,
			Contact2:    contact2,
			Established: parseTime(row[schema.Established]),
		}
		relations = append(relations, relation)
	}
	return relations
}

func populateFriendModels(rows typ.Rows) []model.Friend {

	var friendMs []model.Friend
	for _, row := range rows {
		friendM := model.Friend{
			UserAId:     parseUserId(row[schema.UserAId]),
			UserBId:     parseUserId(row[schema.UserBId]),
			FriendSince: parseTime(row[schema.FriendSince]),
		}
		friendMs = append(friendMs, friendM)
	}

	return friendMs
}

func populateUserModels(rows typ.Rows) []model.User {

	var usrMs []model.User
	for _, row := range rows {
		usrM := model.User{
			Id:      parseUserId(row[schema.UserId]),
			Name:    parseString(row[schema.Name]),
			Email:   parseEmail(row[schema.Email]),
			PwdHash: parsePwdHash(row[schema.PwdHash]),
			Joined:  parseTime(row[schema.CreatedAt]),
		}
		usrMs = append(usrMs, usrM)
	}

	return usrMs
}

func populateMemberModels(rows typ.Rows) []model.Member {

	var mbrMs []model.Member
	for _, row := range rows {
		mbrM := model.Member{
			ChatId: parseChatId(row[schema.ChatId]),
			UserId: parseUserId(row[schema.UserId]),
		}
		mbrMs = append(mbrMs, mbrM)
	}

	return mbrMs
}

func parseTime(v any) time.Time {
	joined, ok := v.(time.Time)
	if !ok {
		log.Fatalf("parseTime: v does not hold a MyType (it is %T)", v)
	}
	return joined
}

func parseEmail(v any) cred.Email {
	email, ok := v.(string)
	if !ok {
		log.Fatalf("parseEmail: v does not hold a MyType (it is %T)", v)
	}
	return cred.Email(email)
}

func parsePwdHash(v any) cred.PwdHash {
	hash, ok := v.([]uint8)
	if !ok {
		log.Fatalf("parsePwdHash: v does not hold a MyType (it is %T)", v)
	}
	return cred.PwdHash(hash)
}

func parseString(v any) string {
	name, ok := v.(string)
	if !ok {
		log.Fatalf("parseName: v does not hold a MyType (it is %T)", v)
	}
	return name
}

func parseMessageId(v any) typ.MessageId {
	mid, ok := v.(int64)
	if !ok {
		log.Fatalf("parseUserId: v does not hold a MyType (it is %T)", v)
	}
	return typ.MessageId(mid)
}

func parseChatId(v any) typ.ChatId {
	cid, ok := v.(int64)
	if !ok {
		log.Fatalf("parseUserId: v does not hold a MyType (it is %T)", v)
	}
	return typ.ChatId(cid)
}

func parseUserId(v any) typ.UserId {
	uid, ok := v.(int64)
	if !ok {
		log.Fatalf("parseUserId: v does not hold a MyType (it is %T)", v)
	}
	return typ.UserId(uid)
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
