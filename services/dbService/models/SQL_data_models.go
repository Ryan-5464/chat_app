package dbServiceModels

import (
	typ "server/types"
	"time"
)

type User struct {
	Id      typ.UserId
	Name    string
	Email   typ.Email
	PwdHash typ.PwdHash
	Joined  time.Time
}

type Chat struct {
	Id        typ.ChatId
	Name      string
	Admin     typ.UserId
	CreatedAt time.Time
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
