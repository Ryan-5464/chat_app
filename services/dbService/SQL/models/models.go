package models

import (
	cred "server/services/authService/credentials"
	typ "server/types"
	"time"
)

type User struct {
	Id      typ.UserId
	Name    string
	Email   cred.Email
	PwdHash cred.PwdHash
	Joined  time.Time
}

type Chat struct {
	Id        typ.ChatId
	Name      string
	AdminId   typ.UserId
	CreatedAt time.Time
}

type ContactChat struct {
	Id            typ.ChatId
	Member1Id     typ.UserId
	Member2Id     typ.UserId
	CreatedAt     time.Time
	LastMessageAt time.Time
}

type Message struct {
	Id         typ.MessageId
	UserId     typ.UserId
	ChatId     typ.ChatId
	ReplyId    typ.MessageId
	Text       string
	CreatedAt  time.Time
	LastEditAt time.Time
}

type Member struct {
	UserId        typ.UserId
	ChatId        typ.ChatId
	LastReadMsgId typ.MessageId
	Joined        time.Time
}
