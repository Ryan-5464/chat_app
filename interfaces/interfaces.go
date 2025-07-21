package interfaces

import (
	"server/data/entities"
	sess "server/services/authService/session"
)

type UserRepository interface {
	NewUser()
	EditUser()
	DeleteUser()
	GetUsers()
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
	NewMessage()
	EditMessage()
	DeleteMessage()
	GetMessages()
	GetMessage()
	CountUnreadMessages()
}

type DbService interface {
	GetUsers() []entities.User
	GetChats() []entities.Chat
	GetMessages() []entities.Message
}

type AuthService interface {
	ValidateAndRefreshSession(JWEtoken string) (sess.Session, error)
}

type ChatService interface {
	GetChats() ([]entities.Chat, error)
}

type MessageService interface {
	GetMessages() ([]entities.Message, error)
}
