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

func (dbs *DbService) GetUsers(userIds []typ.UserId) ([]model.User, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("userIds %v", userIds))

	ids := ToAnySlice(userIds)
	query := selectAllFromWhereIn(schema.UserTable, schema.UserId, ids)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, ids...)
	if err != nil {
		return []model.User{}, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateUserModels(rows), nil
}

func (dbs *DbService) GetMembers(chatIds []typ.ChatId) ([]model.Member, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatIds: %v", chatIds))

	ids := ToAnySlice(chatIds)
	query := selectAllFromWhereIn(schema.MemberTable, schema.ChatId, ids...)

	dbs.lgr.DLog(fmt.Sprintf("query: %v", query))

	rows, err := dbs.db.Read(query, ids...)
	if err != nil {
		return []model.Member{}, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows: %v", rows))

	return populateMemberModels(rows), nil
}

func (dbs *DbService) GetChatUsers(chatId typ.ChatId) ([]model.User, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatId: %v", chatId))

	var chatIds = []typ.ChatId{chatId}
	members, err := dbs.GetMembers(chatIds)
	if err != nil {
		return []model.User{}, err
	}

	dbs.lgr.DLog(fmt.Sprintf("members: %v", members))

	if len(members) == 0 {
		return []model.User{}, errors.New("members missing for chat")
	}

	var userIds []typ.UserId
	for _, member := range members {
		userIds = append(userIds, member.UserId)
	}

	dbs.lgr.DLog(fmt.Sprintf("userIds: %v", userIds))

	if len(userIds) == 0 {
		return []model.User{}, err
	}

	return dbs.GetUsers(userIds)
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

func (dbs *DbService) NewChat(chatName string, adminId typ.UserId) (*model.Chat, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatname %s: adminId %v", chatName, adminId))

	query := insertIntoValues(schema.ChatTable, schema.Name, schema.AdminId)

	dbs.lgr.DLog(fmt.Sprintf("query %s", query))

	res, err := dbs.db.Create(query, chatName, adminId)
	if err != nil {
		return nil, err
	}

	chatId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("chatId %v", chatId))

	return dbs.GetChat(typ.ChatId(chatId))
}

func (dbs *DbService) NewMember(chatId typ.ChatId, userId typ.UserId) error {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("chatId %v: userId %v", chatId, userId))

	query := insertIntoValues(schema.MemberTable, schema.UserId, schema.ChatId)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	_, err := dbs.db.Create(query, chatId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (dbs *DbService) CreateUser(userName string, userEmail cred.Email, pwdHash cred.PwdHash) (typ.LastInsertId, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("values: name - %v email - %v pwdHash - %v", userName, userEmail, pwdHash))

	query := insertIntoValues(schema.UserTable, schema.Name, schema.Email, schema.PwdHash)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	res, err := dbs.db.Create(query, userName, userEmail, pwdHash)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	dbs.lgr.DLog(fmt.Sprintf("lastInsertId %v", lastInsertId))

	return typ.LastInsertId(lastInsertId), nil
}

func (dbs *DbService) GetNewUser(lastInsertId typ.LastInsertId) (*model.User, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("lastInsertId %v", lastInsertId))

	query := selectAllFromWhereEqualTo(schema.UserTable, schema.UserId)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	rows, err := dbs.db.Read(query, lastInsertId)
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows %v", rows))

	return populateUserModel(rows), nil
}

func (dbs *DbService) NewUser(u model.User) (*model.User, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("user: %v", u))

	query := insertIntoValues(schema.UserTable, schema.Name, schema.Email, schema.PwdHash)

	dbs.lgr.DLog(fmt.Sprintf("query: %v", query))

	res, err := dbs.db.Create(query, u.Name, u.Email, u.PwdHash)
	if err != nil {
		return nil, err
	}

	newUsrId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("newUsrId: %v", newUsrId))

	return dbs.GetUser(typ.UserId(newUsrId))
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

func (dbs *DbService) NewMessage(m model.Message) (*model.Message, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("newmessage %v", m))

	query := insertIntoValues(schema.MessageTable, schema.UserId, schema.ChatId, schema.ReplyId, schema.MsgText)

	dbs.lgr.DLog(fmt.Sprintf("query %v", query))

	res, err := dbs.db.Create(query, m.UserId, m.ChatId, m.ReplyId, m.Text)
	if err != nil {
		return nil, err
	}

	newMsgId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("newmessageId %v", newMsgId))

	return dbs.GetMessage(typ.MessageId(newMsgId))
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

func (dbs *DbService) NewContactChat(member1Id typ.UserId, member2Id typ.UserId) (*model.ContactChat, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("member1Id: %v - member2Id: %s", member1Id, member2Id))

	qb := qbuilder.NewQueryBuilder()

	contactChatTable := qb.Table(schema.ContactChatTable)
	member1IdF := qb.Field(schema.Member1Id)
	member2IdF := qb.Field(schema.Member2Id)

	query := qb.INSERT_INTO(contactChatTable, member1IdF, member2IdF).
		VALUES(member1Id, member2Id).Build()

	dbs.lgr.DLog(fmt.Sprintf("query %s", query))

	values := []any{member1Id, member2Id}
	rows, err := dbs.db.Read(query, values...)
	if err != nil {
		return nil, err
	}

	dbs.lgr.DLog(fmt.Sprintf("rows: %v", rows))

	return populateContactChatModel(rows), nil
}

func (dbs *DbService) GetContactChats(userId typ.UserId) ([]model.ContactChat, error) {
	dbs.lgr.LogFunctionInfo()

	dbs.lgr.DLog(fmt.Sprintf("userId: %v", userId))

	qb := qbuilder.NewQueryBuilder()

	contactChatTable := qb.Table(schema.ContactChatTable)
	member1IdF := qb.Field(schema.Member1Id)
	member2IdF := qb.Field(schema.Member2Id)

	query := qb.SELECT(qb.All()).FROM(contactChatTable).
		WHERE(member1IdF, qb.EqualTo()).OR(member2IdF, qb.EqualTo()).Build()

	dbs.lgr.DLog(fmt.Sprintf("query: %s", query))

	rows, err := dbs.db.Read(query, userId, userId)
	if err != nil {
		return []model.ContactChat{}, nil
	}

	dbs.lgr.DLog(fmt.Sprintf("rows: %v", rows))

	return populateContactChatModels(rows), nil
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

func populateContactChatModel(rows typ.Rows) *model.ContactChat {
	contactChats := populateContactChatModels(rows)
	if len(contactChats) == 0 {
		return nil
	}
	contactChat := contactChats[0]
	return &contactChat
}

func populateContactChatModels(rows typ.Rows) []model.ContactChat {
	if len(rows) == 0 {
		return []model.ContactChat{}
	}

	var contactChats []model.ContactChat
	for _, row := range rows {
		contactChat := model.ContactChat{
			Id:            parseChatId(row[schema.ChatId]),
			Member1Id:     parseUserId(row[schema.Member1Id]),
			Member2Id:     parseUserId(row[schema.Member2Id]),
			CreatedAt:     parseTime(row[schema.CreatedAt]),
			LastMessageAt: parseTime(row[schema.LastMsgAt]),
		}
		contactChats = append(contactChats, contactChat)
	}
	return contactChats
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
