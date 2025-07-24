package interfaces

import (
	"server/data/entities"
	sess "server/services/authService/session"
	model "server/services/dbService/SQL/models"
	typ "server/types"
)

type UserRepository interface {
	NewUser(usr entities.User) error
	EditUser()
	DeleteUser()
	GetUsers(chatId typ.ChatId) ([]entities.User, error)
	GetUser()
}

type ChatRepository interface {
	NewChat()
	EditChat()
	DeleteChat()
	GetChats()
	GetChat()
	CountMembers()
}

type MessageRepository interface {
	NewMessage(msg entities.Message) (entities.Message, error)
	EditMessage()
	DeleteMessage()
	GetMessages()
	GetMessage()
	CountUnreadMessages()
}

type DbService interface {
	GetUsers(chatId typ.ChatId) ([]model.User, error)
	GetChats() []entities.Chat
	GetMessages(msgIds []typ.MessageId) ([]model.Message, error)
	NewMessage(msgM model.Message) (model.Message, error)
	NewUser(usrM model.User) error
}

type AuthService interface {
	ValidateAndRefreshSession(JWEtoken string) (sess.Session, error)
	NewSession(userId typ.UserId) (sess.Session, error)
}

type ChatService interface {
	GetChats() ([]entities.Chat, error)
}

type MessageService interface {
	GetMessages(chatId typ.ChatId) ([]entities.Message, error)
}
