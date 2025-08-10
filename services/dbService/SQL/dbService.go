package SQL

import (
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

func (dbs *DbService) UpdateUserEmail(newEmail cred.Email, userId typ.UserId) error {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("newEmail %v", newEmail))

	query := updateWhereEqualTo(schema.UserTable, schema.UserId, schema.Email)

	return dbs.db.Update(query, newEmail, userId)
}

func (dbs *DbService) GetChat(chatId typ.ChatId) (*model.Chat, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatId %v", chatId))

	query := selectAllFromWhereEqualTo(schema.ChatTable, schema.ChatId)

	dbs.lgr.DLog(fmt.Sprintf("query %s", query))

	rows, err := dbs.db.Read(query, chatId)
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateChatModel(rows), nil
}

func (dbs *DbService) GetChats(chatIds []typ.ChatId) ([]model.Chat, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatIds %v", chatIds))

	ids := ToAnySlice(chatIds)

	query := selectAllFromWhereIn(schema.ChatTable, schema.ChatId, ids...)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, ids...)
	if err != nil {
		return []model.Chat{}, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateChatModels(rows), nil
}

func (dbs *DbService) GetUserChats(userId typ.UserId) ([]model.Chat, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("userId %v", userId))

	query := selectAllFromWhereEqualTo(schema.ChatTable, schema.AdminId)

	rows, err := dbs.db.Read(query, userId)
	if err != nil {
		return []model.Chat{}, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateChatModels(rows), nil
}

func (dbs *DbService) GetUser(userId typ.UserId) (*model.User, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("userId %v", userId))

	query := selectAllFromWhereEqualTo(schema.UserTable, schema.UserId)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, userId)
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateUserModel(rows), nil
}

func (dbs *DbService) GetUsers(userIds []typ.UserId) ([]model.User, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("userIds %v", userIds))

	ids := ToAnySlice(userIds)
	query := selectAllFromWhereIn(schema.UserTable, schema.UserId, ids...)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, ids...)
	if err != nil {
		return []model.User{}, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateUserModels(rows), nil
}

func (dbs *DbService) GetMembers(chatId typ.ChatId) ([]model.Member, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatId: %v", chatId))

	query := selectAllFromWhereIn(schema.MemberTable, schema.ChatId, chatId)

	dbs.lgr.DLog(fmt.Sprintf("query: %v", query))

	rows, err := dbs.db.Read(query, chatId)
	if err != nil {
		return []model.Member{}, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows: %v", rows))

	return populateMemberModels(rows), nil
}

func (dbs *DbService) GetChatMessages(chatId typ.ChatId) ([]model.Message, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatId: %v", chatId))

	query := selectAllFromWhereEqualTo(schema.MessageTable, schema.ChatId)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, chatId)
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateMessageModels(rows), nil
}

func (dbs *DbService) GetMessages(messageIds []typ.MessageId) ([]model.Message, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("messageIds: %v", messageIds))

	ids := ToAnySlice(messageIds)
	query := selectAllFromWhereIn(schema.MessageTable, schema.MessageId, ids...)

	dbs.lgr.DLog(fmt.Sprintf("query: %v", query))

	rows, err := dbs.db.Read(query, ids...)
	if err != nil {
		return []model.Message{}, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows: %v", rows))

	return populateMessageModels(rows), nil
}

func (dbs *DbService) CreateChat(chatName string, adminId typ.UserId) (typ.LastInsertId, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatname %s: adminId %v", chatName, adminId))

	query := insertIntoValues(schema.ChatTable, schema.Name, schema.AdminId)

	dbs.lgr.DLog(fmt.Sprintf("query %s", query))

	res, err := dbs.db.Create(query, chatName, adminId)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return typ.LastInsertId(lastInsertId), nil
}

func (dbs *DbService) CreateMember(chatId typ.ChatId, userId typ.UserId) error {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatId %v: userId %v", chatId, userId))

	query := insertIntoValues(schema.MemberTable, schema.ChatId, schema.UserId)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	_, err := dbs.db.Create(query, chatId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DbService) CreateUser(userName string, email cred.Email, pwdHash cred.PwdHash) (typ.LastInsertId, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("user %s: email %v: pwdHash %v", userName, email, pwdHash))

	query := insertIntoValues(schema.UserTable, schema.Name, schema.Email, schema.PwdHash)

	dbs.lgr.DLog(fmt.Sprintf("query: %v", query))

	res, err := dbs.db.Create(query, userName, email, pwdHash)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return typ.LastInsertId(lastInsertId), nil
}

func (dbs *DbService) CreateMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (typ.LastInsertId, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("userId %v: chatId %v: replyId %v: text %s", userId, chatId, replyId, text))

	query := insertIntoValues(schema.MessageTable, schema.UserId, schema.ChatId, schema.ReplyId, schema.MsgText)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	res, err := dbs.db.Create(query, userId, chatId, replyId, text)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return typ.LastInsertId(lastInsertId), nil
}

func (dbs *DbService) CreateContactMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (typ.LastInsertId, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("userId %v: chatId %v: replyId %v: text %s", userId, chatId, replyId, text))

	query := insertIntoValues(schema.ContactMessageTable, schema.UserId, schema.ChatId, schema.ReplyId, schema.MsgText)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	res, err := dbs.db.Create(query, userId, chatId, replyId, text)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return typ.LastInsertId(lastInsertId), nil
}

func (dbs *DbService) GetContactMessages(chatId typ.ChatId) ([]model.Message, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatId %v", chatId))

	query := selectAllFromWhereEqualTo(schema.ContactMessageTable, schema.ChatId)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, chatId)
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateMessageModels(rows), nil

}

func (dbs *DbService) GetContactMessage(messageId typ.MessageId) (*model.Message, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("messageId %v", messageId))

	query := selectAllFromWhereEqualTo(schema.ContactMessageTable, schema.MessageId)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, messageId)
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateMessageModel(rows), nil
}

func (dbs *DbService) GetMessage(messageId typ.MessageId) (*model.Message, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("messageId %v", messageId))

	query := selectAllFromWhereEqualTo(schema.MessageTable, schema.MessageId)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, messageId)
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateMessageModel(rows), nil
}

func (dbs *DbService) FindUser(email cred.Email) (*model.User, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("email %v", email))

	query := selectAllFromWhereEqualTo(schema.UserTable, schema.Email)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, email)
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateUserModel(rows), nil
}

func (dbs *DbService) FindUsers(e []cred.Email) ([]model.User, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("emails %v", e))

	emails := ToAnySlice(e)
	query := selectAllFromWhereIn(schema.UserTable, schema.Email, emails...)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, emails...)
	if err != nil {
		return []model.User{}, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateUserModels(rows), nil
}

func (dbs *DbService) CreateContact(id1 typ.UserId, id2 typ.ContactId) (typ.LastInsertId, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("id1 %v: id2 %v", id1, id2))

	query := insertIntoValues(schema.ContactTable, schema.Id1, schema.Id2)

	dbs.lgr.DLog(fmt.Sprintf("query %s", query))

	res, err := dbs.db.Create(query, id1, id2)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return typ.LastInsertId(lastInsertId), nil
}

func (dbs *DbService) GetContact(chatId typ.ChatId) (*model.Contact, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatId: %v", chatId))

	query := selectAllFromWhereEqualTo(schema.ContactTable, schema.ChatId)

	dbs.lgr.DLog(fmt.Sprintf("query: %s", query))

	rows, err := dbs.db.Read(query, chatId)
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows: %v", rows))

	return populateContactModel(rows), nil
}

func (dbs *DbService) GetContacts(userId typ.UserId) ([]model.Contact, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("userId: %v", userId))

	qb := qbuilder.NewQueryBuilder()

	contactTable := qb.Table(schema.ContactTable)
	id1F := qb.Field(schema.Id1)
	id2F := qb.Field(schema.Id2)

	query := qb.SELECT(qb.All()).FROM(contactTable).
		WHERE(id1F, qb.EqualTo()).OR(id2F, qb.EqualTo()).Build()

	dbs.lgr.DLog(fmt.Sprintf("query: %s", query))

	rows, err := dbs.db.Read(query, userId, userId)
	if err != nil {
		return []model.Contact{}, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows: %v", rows))

	return populateContactModels(rows), nil
}

func selectAllFromWhereEqualTo(table string, field string) string {
	qb := qbuilder.NewQueryBuilder()

	t := qb.Table(table)
	f := qb.Field(field)

	return qb.SELECT(qb.All()).FROM(t).WHERE(f, qb.EqualTo()).Build()
}

func selectAllFromWhereIn(table string, field string, values ...any) string {
	qb := qbuilder.NewQueryBuilder()

	t := qb.Table(table)
	f := qb.Field(field)

	return qb.SELECT(qb.All()).FROM(t).WHERE(f, qb.IN(values...)).Build()
}

func insertIntoValues(table string, fields ...string) string {
	qb := qbuilder.NewQueryBuilder()

	t := qb.Table(table)

	var fs []qbuilder.Field
	for _, field := range fields {
		f := qb.Field(field)
		fs = append(fs, f)
	}

	vals := ToAnySlice(fs)

	return qb.INSERT_INTO(t, fs...).VALUES(vals...).Build()
}

func populateContactModel(rows typ.Rows) *model.Contact {
	contacts := populateContactModels(rows)
	if len(contacts) == 0 {
		return nil
	}
	contact := contacts[0]
	return &contact
}

func populateContactModels(rows typ.Rows) []model.Contact {
	if len(rows) == 0 {
		return []model.Contact{}
	}

	var contact []model.Contact
	for _, row := range rows {
		chat := model.Contact{
			ChatId:        parseChatId(row[schema.ChatId]),
			Id1:           parseUserId(row[schema.Id1]),
			Id2:           parseUserId(row[schema.Id2]),
			CreatedAt:     parseTime(row[schema.CreatedAt]),
			LastMessageAt: parseTime(row[schema.LastMsgAt]),
		}
		contact = append(contact, chat)
	}
	return contact
}

func populateChatModel(rows typ.Rows) *model.Chat {
	chats := populateChatModels(rows)
	if len(chats) == 0 {
		return nil
	}
	chat := chats[0]
	return &chat
}

func populateChatModels(rows typ.Rows) []model.Chat {
	if len(rows) == 0 {
		return []model.Chat{}
	}

	var chatModels []model.Chat
	for _, row := range rows {
		chatModel := model.Chat{
			Id:        parseChatId(row[schema.ChatId]),
			Name:      parseString(row[schema.Name]),
			AdminId:   parseUserId(row[schema.AdminId]),
			CreatedAt: parseTime(row[schema.CreatedAt]),
		}
		chatModels = append(chatModels, chatModel)
	}
	return chatModels
}

func populateMessageModel(rows typ.Rows) *model.Message {
	messages := populateMessageModels(rows)
	if len(messages) == 0 {
		return nil
	}
	message := messages[0]
	return &message
}

func populateMessageModels(rows typ.Rows) []model.Message {
	if len(rows) == 0 {
		return []model.Message{}
	}

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

func populateUserModel(rows typ.Rows) *model.User {
	users := populateUserModels(rows)
	if len(users) == 0 {
		return nil
	}
	user := users[0]
	return &user
}

func populateUserModels(rows typ.Rows) []model.User {
	if len(rows) == 0 {
		return []model.User{}
	}

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

func populateMemberModel(rows typ.Rows) *model.Member {
	members := populateMemberModels(rows)
	if len(members) == 0 {
		return nil
	}
	member := members[0]
	return &member
}

func populateMemberModels(rows typ.Rows) []model.Member {
	if len(rows) == 0 {
		return []model.Member{}
	}

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

func updateWhereEqualTo(table string, conditionField string, updateFields ...string) string {
	qb := qbuilder.NewQueryBuilder()

	t := qb.Table(table)
	cond := qb.Field(conditionField)

	var setFields []qbuilder.Set
	for _, field := range updateFields {
		setFields = append(setFields, qb.SET(qb.Field(field)))
	}

	return qb.UPDATE(t, setFields...).WHERE(cond, qb.EqualTo()).Build()
}
