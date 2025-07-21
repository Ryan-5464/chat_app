package interfaces

import (
	"server/data/entities"
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
