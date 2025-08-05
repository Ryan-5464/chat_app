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
	NewUser(usr ent.User) (*ent.User, error)
	GetChatUsers(chatId typ.ChatId) ([]ent.User, error)
	GetUsers(userIds []typ.UserId) ([]ent.User, error)
	FindUsers(emails []cred.Email) ([]ent.User, error)
	GetContacts(userId typ.UserId) ([]ent.Contact, error)
	AddContact(contact ent.Contact, userId typ.UserId) (*ent.Contact, error)
}

type ChatRepository interface {
	NewChat(chatName string, adminId typ.UserId, chatType typ.ChatType) (*ent.Chat, error)
	GetChats(userId typ.UserId) ([]ent.Chat, error)
	NewMember(chatId typ.ChatId, userId typ.UserId) error
}

type MessageRepository interface {
	NewMessage(msg ent.Message) (*ent.Message, error)
	GetChatMessages(chatId typ.ChatId) ([]ent.Message, error)
}

type DbService interface {
	GetUsers(usrIds []typ.UserId) ([]model.User, error)
	GetChatUsers(chatId typ.ChatId) ([]model.User, error)
	NewUser(userModel model.User) (*model.User, error)
	FindUsers(emails []cred.Email) ([]model.User, error)
	GetContactRelations(userId typ.UserId) ([]model.ContactRelation, error)
	AddContactRelation(userId typ.UserId, contactId typ.UserId) (*model.ContactRelation, error)
	NewMember(chatId typ.ChatId, userId typ.UserId) error
	NewChat(chatName string, adminId typ.UserId, chatType typ.ChatType) (*model.Chat, error)
	GetChats(chatId []typ.ChatId) ([]model.Chat, error)
	GetUserChats(userId typ.UserId) ([]model.Chat, error)
	GetMessages(msgIds []typ.MessageId) ([]model.Message, error)
	GetChatMessages(chatId typ.ChatId) ([]model.Message, error)
	NewMessage(msgM model.Message) (*model.Message, error)
	GetPrivateChatIdsForContacts(userId typ.UserId) ([]model.Member, error)
	Close()
}

type AuthService interface {
	ValidateAndRefreshSession(JWEtoken string) (sess.Session, error)
	NewSession(userId typ.UserId) (sess.Session, error)
}

type ChatService interface {
	GetChats(userId typ.UserId) ([]ent.Chat, error)
	NewChat(chat ent.Chat) (*ent.Chat, error)
}

type MessageService interface {
	GetChatMessages(chatId typ.ChatId) ([]ent.Message, error)
	HandleNewMessage(msg ent.Message) error
	NewMessage(msg ent.Message) (*ent.Message, error)
}

type UserService interface {
	GetUsers(userId []typ.UserId) ([]ent.User, error)
	GetChatUsers(chatId typ.ChatId) ([]ent.User, error)
	NewUser(newUser dto.NewUserInput) (*ent.User, error)
	FindUsers(emails []cred.Email) ([]ent.User, error)
	AddContact(a dto.AddContactInput) (*ent.Contact, error)
	GetContacts(userId typ.UserId) ([]ent.Contact, error)
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
