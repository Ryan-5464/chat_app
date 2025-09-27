package interfaces

import (
	ent "server/data/entities"
	cred "server/services/auth/credentials"
	sess "server/services/auth/session"
	model "server/services/db/SQL/models"
	typ "server/types"
)

type UserRepository interface {
	AddContact(ct typ.ContactId, contactname string, ce cred.Email, u typ.UserId) (*ent.Contact, error)
	FindUser(e cred.Email) (*model.User, error)
	FindUsers(es []cred.Email) ([]ent.User, error)
	GetChatUsers(c typ.ChatId) ([]ent.User, error)
	GetContact(c typ.ChatId, u typ.UserId) (*ent.Contact, error)
	GetContacts(u typ.UserId) ([]ent.Contact, error)
	GetUser(u typ.UserId) (*ent.User, error)
	GetUsers(us []typ.UserId) ([]ent.User, error)
	NewUser(username string, e cred.Email, p cred.PwdHash) (*ent.User, error)
	RemoveContact(ct typ.ContactId, u typ.UserId) error
	EditUserName(username string, u typ.UserId) error
	GetUserByEmail(e cred.Email) (*ent.User, error)
}

type ChatRepository interface {
	DeleteChat(c typ.ChatId) error
	GetChat(c typ.ChatId) (*ent.Chat, error)
	GetChats(u typ.UserId) ([]ent.Chat, error)
	NewChat(chatname string, adminId typ.UserId) (*ent.Chat, error)
	NewMember(c typ.ChatId, u typ.UserId) error
	GetMember(c typ.ChatId, u typ.UserId) (*ent.Member, error)
	GetMembers(c typ.ChatId) ([]ent.Member, error)
	RemoveChatMember(c typ.ChatId, u typ.UserId) error
	NewChatAdmin(c typ.ChatId, newAdminId typ.UserId) error
	VerifyChatAdmin(c typ.ChatId, u typ.UserId) (bool, error)
	EditChatName(newname string, c typ.ChatId) error
	GetChatMemberships(u typ.UserId) ([]ent.Member, error)
	GetUnreadMessageCount(lastReadMsgId typ.MessageId, c typ.ChatId) (int64, error)
}

type MessageRepository interface {
	GetChatMessages(c typ.ChatId) ([]ent.Message, error)
	GetContactMessages(c typ.ChatId) ([]ent.Message, error)
	NewContactMessage(u typ.UserId, c typ.ChatId, replyId *typ.MessageId, text string) (*ent.Message, error)
	NewMessage(u typ.UserId, c typ.ChatId, replyId *typ.MessageId, text string) (*ent.Message, error)
	DeleteMessage(m typ.MessageId) error
	GetMessage(m typ.MessageId) (*ent.Message, error)
	EditMessage(msgText string, m typ.MessageId) error
	UpdateLastReadMsgId(lastReadMsgId typ.MessageId, c typ.ChatId, u typ.UserId) error
	GetLatestChatMessageId(c typ.ChatId) (typ.MessageId, error)
	GetLatestMessageId() (typ.MessageId, error)
}

type DbService interface {
	FindUser(e cred.Email) (*model.User, error)
	FindUsers(es []cred.Email) ([]model.User, error)

	CreateUser(un string, e cred.Email, p cred.PwdHash) (typ.LastInsertId, error)
	GetUser(us typ.UserId) (*model.User, error)
	GetUsers(us []typ.UserId) ([]model.User, error)
	UpdateUserName(name string, u typ.UserId) error
	GetUserByEmail(e cred.Email) (*model.User, error)

	CreateMember(c typ.ChatId, u typ.UserId) error
	DeleteMember(c typ.ChatId, u typ.UserId) error
	GetMember(c typ.ChatId, u typ.UserId) (*model.Member, error)
	GetMembers(c typ.ChatId) ([]model.Member, error)
	GetMemberships(u typ.UserId) ([]model.Member, error)

	CreateChat(chatname string, adminId typ.UserId) (typ.LastInsertId, error)
	DeleteChat(c typ.ChatId) error
	GetChat(c typ.ChatId) (*model.Chat, error)
	GetChats(c []typ.ChatId) ([]model.Chat, error)
	GetUserChats(u typ.UserId) ([]model.Chat, error)
	UpdateChatAdmin(c typ.ChatId, u typ.UserId) error
	UpdateChatName(newname string, c typ.ChatId) error

	CreateContactMessage(u typ.UserId, c typ.ChatId, replyId *typ.MessageId, text string) (typ.LastInsertId, error)
	CreateMessage(u typ.UserId, c typ.ChatId, replyId *typ.MessageId, text string) (typ.LastInsertId, error)
	DeleteMessage(m typ.MessageId) error
	GetChatMessages(c typ.ChatId) ([]model.Message, error)
	GetContactMessage(m typ.MessageId) (*model.Message, error)
	GetContactMessages(c typ.ChatId) ([]model.Message, error)
	GetMessage(m typ.MessageId) (*model.Message, error)
	GetMessages(msgIds []typ.MessageId) ([]model.Message, error)
	UpdateMessage(msgtext string, m typ.MessageId) error
	UpdateLastReadMsgId(lastReadMsgId typ.MessageId, c typ.ChatId, u typ.UserId) error
	GetUnreadMessageCount(lastReadMsgId typ.MessageId, c typ.ChatId) (int64, error)
	GetLatestChatMessageId(c typ.ChatId) (typ.MessageId, error)
	GetLatestMessageId() (typ.MessageId, error)

	CreateContact(id1 typ.UserId, id2 typ.ContactId) (typ.LastInsertId, error)
	GetContact(c typ.ChatId) (*model.Contact, error)
	GetContacts(u typ.UserId) ([]model.Contact, error)
	DeleteContact(ct typ.ContactId, u typ.UserId) error

	Close()
}

type AuthService interface {
	NewSession(u typ.UserId) (sess.Session, error)
	ValidateAndRefreshSession(JWEtoken string) (sess.Session, error)
}

type ChatService interface {
	GetChats(u typ.UserId) ([]ent.Chat, error)
	NewChat(chatname string, adminId typ.UserId) (*ent.Chat, error)
	LeaveChat(c typ.ChatId, u typ.UserId) ([]ent.Chat, error)
	EditChatName(newname string, c typ.ChatId, u typ.UserId) error
	GetChatMembers(c typ.ChatId) ([]ent.Member, error)
	AddMember(e cred.Email, c typ.ChatId) (typ.UserId, error)
	GetChatMember(c typ.ChatId, u typ.UserId) (*ent.Member, error)
	RemoveMember(c typ.ChatId, u typ.UserId, adminId typ.UserId) error
	GetUnreadMessageCount(c typ.ChatId, u typ.UserId) (int64, error)
}

type MessageService interface {
	GetChatMessages(c typ.ChatId, u typ.UserId) ([]ent.Message, error)
	GetContactMessages(c typ.ChatId, u typ.UserId) ([]ent.Message, error)
	HandleNewContactMessage(u typ.UserId, c typ.ChatId, replyId *typ.MessageId, msgText string) error
	HandleNewMessage(u typ.UserId, c typ.ChatId, replyId *typ.MessageId, msgText string) error
	DeleteMessage(m typ.MessageId) error
	EditMessage(msgText string, m typ.MessageId) (*ent.Message, error)
	UpdateLastReadMsgId(lastReadMsgId typ.MessageId, c typ.ChatId, u typ.UserId) error
	GetLatestChatMessageId(c typ.ChatId) (typ.MessageId, error)
	GetLatestMessageId() (typ.MessageId, error)
}

type UserService interface {
	AddContact(e cred.Email, u typ.UserId) (*ent.Contact, error)
	FindUsers(es []cred.Email) ([]ent.User, error)
	GetChatUsers(c typ.ChatId) ([]ent.User, error)
	GetContact(c typ.ChatId, u typ.UserId) (*ent.Contact, error)
	GetContacts(u typ.UserId) ([]ent.Contact, error)
	GetUser(u typ.UserId) (*ent.User, error)
	GetUsers(u []typ.UserId) ([]ent.User, error)
	NewUser(name string, e cred.Email, p cred.PwdHash) (*ent.User, error)
	RemoveContact(ct typ.ContactId, u typ.UserId) error
	EditUserName(name string, u typ.UserId) error
	GetUserByEmail(e cred.Email) (*ent.User, error)
}

type ConnectionService interface {
	DisconnectUser(u typ.UserId)
	GetConnection(u typ.UserId) Socket
	StoreConnection(conn Socket, u typ.UserId)
	GetActiveConnections() map[typ.UserId]Socket
	ChangeOnlineStatus(status string, u typ.UserId) error
	GetOnlineStatus(userId typ.UserId) string
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
