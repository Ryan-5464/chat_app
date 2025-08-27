package entities

import (
	cred "server/services/auth/credentials"
	tkn "server/services/auth/jwetoken"
	typ "server/types"
	"time"
)

type Message struct {
	Id            typ.MessageId
	UserId        typ.UserId
	ChatId        typ.ChatId
	ReplyId       typ.MessageId
	Author        string
	Text          string
	CreatedAt     time.Time
	LastEditAt    time.Time
	IsUserMessage bool
}

type Chat struct {
	Id                 typ.ChatId
	Name               string
	AdminId            typ.UserId
	AdminName          string
	MemberCount        int64
	UnreadMessageCount int64
	CreatedAt          time.Time
	UserIsAdmin        bool
}

type User struct {
	Id      typ.UserId
	Name    string
	Email   cred.Email
	PwdHash cred.PwdHash
	Joined  time.Time
	Token   tkn.JWE
}

type Contact struct {
	Id            typ.ContactId
	Name          string
	Email         cred.Email
	KnownSince    time.Time
	OnlineStatus  string
	ContactChatId typ.ChatId
}

type Member struct {
	ChatId        typ.ChatId
	UserId        typ.UserId
	LastReadMsgId typ.MessageId
	Joined        time.Time
	Name          string
	Email         cred.Email
}
