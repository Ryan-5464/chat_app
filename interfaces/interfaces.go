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
	GetUsers(userIds []typ.UserId) ([]ent.User, error)
	GetUser(userId typ.UserId) (ent.User, error)
}

type ChatRepository interface {
	NewChat(chat ent.Chat) (ent.Chat, error)
	// EditChat()
	// DeleteChat()
	GetChats()
	// GetChat()
	// CountMembers()
}

type MessageRepository interface {
	NewMessage(msg ent.Message) (ent.Message, error)
	// EditMessage()
	// DeleteMessage()
	// GetMessages()
	// GetMessage()
	// CountUnreadMessages()
}

type DbService interface {
	FindUser(email cred.Email) (model.User, error)
	GetUser(usrId typ.UserId) (model.User, error)
	GetUsers(usrIds []typ.UserId) ([]model.User, error)
	NewChat(chat model.Chat) (model.Chat, error)
	GetChats(chatIds []typ.ChatId) ([]model.Chat, error)
	GetMessage(msgId typ.MessageId) (model.Message, error)
	GetMessages(msgIds []typ.MessageId) ([]model.Message, error)
	NewMessage(msgM model.Message) (model.Message, error)
	NewUser(usrM model.User) (model.User, error)
	Close()
}

type AuthService interface {
	ValidateAndRefreshSession(JWEtoken string) (sess.Session, error)
	NewSession(userId typ.UserId) (sess.Session, error)
}

type ChatService interface {
	GetChats() ([]ent.Chat, error)
	NewChat(chat ent.Chat) (ent.Chat, error)
}

type MessageService interface {
	GetMessages(chatId typ.ChatId) ([]ent.Message, error)
	HandleNewMessage(msg ent.Message) error
}

type UserService interface {
	FindUser(usr ent.User) (ent.User, error)
	GetUsers(chatId typ.ChatId) ([]ent.User, error)
	GetUser(userId typ.UserId) (ent.User, error)
	NewUser(user ent.User) (ent.User, error)
}

type ConnectionService interface {
	StoreConnection(conn Socket, userId typ.UserId)
	GetConnection(userId typ.UserId) Socket
}

type Socket interface {
	ReadJSON(v any) error
	WriteJSON(v any) error
	Close() error
}

type Logger interface {
	Log(message string)
	LogError(err error)
	LogFunctionInfo()
}
