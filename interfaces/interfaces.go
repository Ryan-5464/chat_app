package interfaces

import (
	dto "server/data/DTOs"
	ent "server/data/entities"
	cred "server/services/authService/credentials"
	sess "server/services/authService/session"
	model "server/services/dbService/SQL/models"
	typ "server/types"
)

type UserRepository interface {
	AddContact(contactId typ.ContactId, contactName string, contactEmail cred.Email, userId typ.UserId) (*ent.Contact, error)
	FindUser(email cred.Email) (*model.User, error)
	FindUsers(emails []cred.Email) ([]ent.User, error)
	GetChatUsers(chatId typ.ChatId) ([]ent.User, error)
	GetContact(chatId typ.ChatId, userId typ.UserId) (*ent.Contact, error)
	GetContacts(userId typ.UserId) ([]ent.Contact, error)
	GetUser(userId typ.UserId) (*ent.User, error)
	GetUsers(userIds []typ.UserId) ([]ent.User, error)
	NewUser(userName string, userEmail cred.Email, pwdHash cred.PwdHash) (*ent.User, error)
}

type ChatRepository interface {
	GetChats(userId typ.UserId) ([]ent.Chat, error)
	NewChat(chatName string, adminId typ.UserId) (*ent.Chat, error)
	NewMember(chatId typ.ChatId, userId typ.UserId) error
}

type MessageRepository interface {
	GetChatMessages(chatId typ.ChatId) ([]ent.Message, error)
	NewContactMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (*ent.Message, error)
	NewMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (*ent.Message, error)
}

type DbService interface {
	FindUser(email cred.Email) (*model.User, error)
	FindUsers(emails []cred.Email) ([]model.User, error)

	CreateUser(userName string, email cred.Email, pwdHash cred.PwdHash) (typ.LastInsertId, error)
	GetUser(usrIds typ.UserId) (*model.User, error)
	GetUsers(usrIds []typ.UserId) ([]model.User, error)

	CreateMember(chatId typ.ChatId, userId typ.UserId) error
	GetMembers(chatId typ.ChatId) ([]model.Member, error)

	CreateChat(chatName string, adminId typ.UserId) (typ.LastInsertId, error)
	GetChat(chatId typ.ChatId) (*model.Chat, error)
	GetChats(chatId []typ.ChatId) ([]model.Chat, error)
	GetUserChats(userId typ.UserId) ([]model.Chat, error)

	CreateContactMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (typ.LastInsertId, error)
	CreateMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (typ.LastInsertId, error)
	GetChatMessages(chatId typ.ChatId) ([]model.Message, error)
	GetContactMessage(messageId typ.MessageId) (*model.Message, error)
	GetMessage(msgId typ.MessageId) (*model.Message, error)
	GetMessages(msgIds []typ.MessageId) ([]model.Message, error)

	CreateContact(id1 typ.UserId, id2 typ.ContactId) (*model.Contact, error)
	GetContact(chatId typ.ChatId) (*model.Contact, error)
	GetContacts(userId typ.UserId) ([]model.Contact, error)

	Close()
}

type AuthService interface {
	NewSession(userId typ.UserId) (sess.Session, error)
	ValidateAndRefreshSession(JWEtoken string) (sess.Session, error)
}

type ChatService interface {
	GetChats(userId typ.UserId) ([]ent.Chat, error)
	NewChat(chatName string, adminId typ.UserId) (*ent.Chat, error)
}

type MessageService interface {
	GetChatMessages(chatId typ.ChatId) ([]ent.Message, error)
	HandleNewContactMessage(mi dto.NewMessageInput) error
	HandleNewMessage(mi dto.NewMessageInput) error
}

type UserService interface {
	AddContact(a dto.AddContactInput) (*ent.Contact, error)
	FindUsers(emails []cred.Email) ([]ent.User, error)
	GetChatUsers(chatId typ.ChatId) ([]ent.User, error)
	GetContact(chatId typ.ChatId, userId typ.UserId) (*ent.Contact, error)
	GetContacts(userId typ.UserId) ([]ent.Contact, error)
	GetUser(userId typ.UserId) (*ent.User, error)
	GetUsers(userId []typ.UserId) ([]ent.User, error)
	NewUser(newUser dto.NewUserInput) (*ent.User, error)
}

type ConnectionService interface {
	DisconnectUser(userId typ.UserId)
	GetConnection(userId typ.UserId) Socket
	StoreConnection(conn Socket, userId typ.UserId)
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
