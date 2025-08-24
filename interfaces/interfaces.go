package interfaces

import (
	dto "server/data/DTOs"
	"server/data/entities"
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
	RemoveContact(contactId typ.ContactId, userId typ.UserId) error
	EditUserName(name string, userId typ.UserId) error
	GetUserByEmail(email cred.Email) (*ent.User, error)
}

type ChatRepository interface {
	DeleteChat(chatId typ.ChatId) error
	GetChat(chatId typ.ChatId) (*ent.Chat, error)
	GetChats(userId typ.UserId) ([]ent.Chat, error)
	NewChat(chatName string, adminId typ.UserId) (*ent.Chat, error)
	NewMember(chatId typ.ChatId, userId typ.UserId) error
	GetMembers(chatId typ.ChatId) ([]ent.Member, error)
	RemoveChatMember(chatId typ.ChatId, userId typ.UserId) error
	NewChatAdmin(chatId typ.ChatId, newAdminId typ.UserId) error
	VerifyChatAdmin(chatId typ.ChatId, userId typ.UserId) (bool, error)
	EditChatName(newName string, chatId typ.ChatId) error
	GetChatMemberships(userId typ.UserId) ([]ent.Member, error)
}

type MessageRepository interface {
	GetChatMessages(chatId typ.ChatId) ([]ent.Message, error)
	GetContactMessages(chatId typ.ChatId) ([]entities.Message, error)
	NewContactMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (*ent.Message, error)
	NewMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (*ent.Message, error)
	DeleteMessage(messageId typ.MessageId) error
	GetMessage(msgId typ.MessageId) (*ent.Message, error)
	EditMessage(msgText string, msgId typ.MessageId) error
}

type DbService interface {
	FindUser(email cred.Email) (*model.User, error)
	FindUsers(emails []cred.Email) ([]model.User, error)

	CreateUser(userName string, email cred.Email, pwdHash cred.PwdHash) (typ.LastInsertId, error)
	GetUser(usrIds typ.UserId) (*model.User, error)
	GetUsers(usrIds []typ.UserId) ([]model.User, error)
	UpdateUserName(name string, userId typ.UserId) error
	GetUserByEmail(email cred.Email) (*model.User, error)

	CreateMember(chatId typ.ChatId, userId typ.UserId) error
	DeleteMember(chatId typ.ChatId, userId typ.UserId) error
	GetMembers(chatId typ.ChatId) ([]model.Member, error)
	GetMemberships(userId typ.UserId) ([]model.Member, error)

	CreateChat(chatName string, adminId typ.UserId) (typ.LastInsertId, error)
	DeleteChat(chatId typ.ChatId) error
	GetChat(chatId typ.ChatId) (*model.Chat, error)
	GetChats(chatId []typ.ChatId) ([]model.Chat, error)
	GetUserChats(userId typ.UserId) ([]model.Chat, error)
	UpdateChatAdmin(chatId typ.ChatId, userId typ.UserId) error
	UpdateChatName(newName string, chatId typ.ChatId) error

	CreateContactMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (typ.LastInsertId, error)
	CreateMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (typ.LastInsertId, error)
	DeleteMessage(messageId typ.MessageId) error
	GetChatMessages(chatId typ.ChatId) ([]model.Message, error)
	GetContactMessage(messageId typ.MessageId) (*model.Message, error)
	GetContactMessages(chatId typ.ChatId) ([]model.Message, error)
	GetMessage(msgId typ.MessageId) (*model.Message, error)
	GetMessages(msgIds []typ.MessageId) ([]model.Message, error)
	UpdateMessage(msgtext string, msgId typ.MessageId) error

	CreateContact(id1 typ.UserId, id2 typ.ContactId) (typ.LastInsertId, error)
	GetContact(chatId typ.ChatId) (*model.Contact, error)
	GetContacts(userId typ.UserId) ([]model.Contact, error)
	DeleteContact(contactId typ.ContactId, userId typ.UserId) error

	Close()
}

type AuthService interface {
	NewSession(userId typ.UserId) (sess.Session, error)
	ValidateAndRefreshSession(JWEtoken string) (sess.Session, error)
}

type ChatService interface {
	GetChats(userId typ.UserId) ([]ent.Chat, error)
	NewChat(chatName string, adminId typ.UserId) (*ent.Chat, error)
	LeaveChat(chatId typ.ChatId, userId typ.UserId) ([]ent.Chat, error)
	EditChatName(newName string, chatId typ.ChatId, userId typ.UserId) error
	GetChatMembers(chatId typ.ChatId) ([]ent.Member, error)
	AddMember(email cred.Email, chatId typ.ChatId) (typ.UserId, error)
	GetChatMember(chatId typ.ChatId, userId typ.UserId) (*ent.Member, error)
	RemoveMember(chatId typ.ChatId, userId typ.UserId, adminId typ.UserId) error
}

type MessageService interface {
	GetChatMessages(chatId typ.ChatId, userId typ.UserId) ([]ent.Message, error)
	GetContactMessages(chatId typ.ChatId, userId typ.UserId) ([]entities.Message, error)
	HandleNewContactMessage(u typ.UserId, c typ.ChatId, replyId *typ.MessageId, msgText string) error
	HandleNewMessage(u typ.UserId, c typ.ChatId, replyId *typ.MessageId, msgText string) error
	DeleteMessage(messageId typ.MessageId) error
	EditMessage(msgText string, msgId typ.MessageId) (*ent.Message, error)
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
	RemoveContact(contactId typ.ContactId, userId typ.UserId) error
	EditUserName(name string, userId typ.UserId) error
	GetUserByEmail(email cred.Email) (*ent.User, error)
}

type ConnectionService interface {
	DisconnectUser(userId typ.UserId)
	GetConnection(userId typ.UserId) Socket
	StoreConnection(conn Socket, userId typ.UserId)
	GetActiveConnections() map[typ.UserId]Socket
}

type Socket interface {
	ReadJSON(v any) error
	WriteJSON(v any) error
	Close() error
}

type Logger interface {
	Dbug(message string)
	Dbugf(message string, values ...any)
	Info(message string)
	Infof(message string, values ...any)
	Error(err error)
	Errorf(message string, err error, values ...any)
	LogFunctionInfo()
}
