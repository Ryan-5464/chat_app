package interfaces

import (
	ent "server/data/entities"
	cred "server/services/authService/credentials"
	sess "server/services/authService/session"
	model "server/services/dbService/SQL/models"
	typ "server/types"
)

type UserRepository interface {
	NewUser(usr ent.User) (ent.User, error)
	FindUser(usr ent.User) (ent.User, error)
	// EditUser()
	// DeleteUser()
	GetUsersForChat(chatId typ.ChatId) ([]ent.User, error)
	GetUsers(userIds []typ.UserId) ([]model.User, error)
	FindUsers(emails []cred.Email) ([]model.User, error)
	GetContactRelations(userId typ.UserId) ([]model.ContactRelation, error)
	AddContact(contact ent.Contact, user ent.User) ([]ent.Contact, error)
}

type ChatRepository interface {
	NewChat(chat ent.Chat) (ent.Chat, error)
	// EditChat()
	// DeleteChat()
	GetChats(userId typ.UserId) ([]ent.Chat, error)
	// GetChat()
	// CountMembers()
}

type MessageRepository interface {
	NewMessage(msg ent.Message) (ent.Message, error)
	// EditMessage()
	// DeleteMessage()
	GetChatMessages(chatId typ.ChatId) ([]ent.Message, error)
	// GetMessage()
	// CountUnreadMessages()
}

type DbService interface {
	// SCRAP THE SINGULAR RETURN FUNCTIONS AND ALWAYS RETURN A SLICE, CAN CHECK LEN SLICE == 0
	FindUsers(emails []cred.Email) ([]model.User, error)
	GetUser(usrId typ.UserId) (model.User, error)
	GetUsers(usrIds []typ.UserId) ([]model.User, error)
	GetUsersForChat(chatId typ.ChatId) ([]model.User, error)
	NewMember(chatId typ.ChatId, userId typ.UserId) error
	NewChat(chat model.Chat) (model.Chat, error)
	GetChat(chatId typ.ChatId) (model.Chat, error)
	GetChats(userId typ.UserId) ([]model.Chat, error)
	GetMessage(msgId typ.MessageId) (model.Message, error)
	GetMessages(msgIds []typ.MessageId) ([]model.Message, error)
	GetChatMessages(chatId typ.ChatId) ([]model.Message, error)
	NewMessage(msgM model.Message) (model.Message, error)
	NewUser(usrM model.User) (model.User, error)
	GetContactRelations(userId typ.UserId) ([]model.ContactRelation, error)
	AddContactRelation(userId typ.UserId, contactId typ.UserId) ([]model.ContactRelation, error)
	Close()
}

type AuthService interface {
	ValidateAndRefreshSession(JWEtoken string) (sess.Session, error)
	NewSession(userId typ.UserId) (sess.Session, error)
}

type ChatService interface {
	GetChats(userId typ.UserId) ([]ent.Chat, error)
	NewChat(chat ent.Chat) (ent.Chat, error)
}

type MessageService interface {
	GetChatMessages(chatId typ.ChatId) ([]ent.Message, error)
	HandleNewMessage(msg ent.Message) error
	NewMessage(msg ent.Message) (ent.Message, error)
}

type UserService interface {
	FindUser(usr ent.User) (ent.User, error)
	GetUsers(chatId typ.ChatId) ([]ent.User, error)
	GetUsersForChat(chatId typ.ChatId) ([]ent.User, error)
	GetUser(userId typ.UserId) (ent.User, error)
	NewUser(user ent.User) (ent.User, error)
	GetContacts(userId typ.UserId) ([]ent.Contact, error)
	AddContact(contact ent.Contact, user ent.User) ([]ent.Contact, error)
}

type ConnectionService interface {
	StoreConnection(conn Socket, userId typ.UserId)
	GetConnection(userId typ.UserId) Socket
	DisconnectUser(userId typ.UserId)
}

type Socket interface {
	ReadJSON(v any) error
	WriteJSON(v any) error
	Close() error
}

type Logger interface {
	DLog(message string)
	Log(message string)
	LogError(err error)
	LogFunctionInfo()
}
