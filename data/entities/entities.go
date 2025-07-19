package dataEntities

import (
	typ "server/types"
	"time"
)

type Message struct {
	Id         typ.MessageId
	UserId     typ.UserId
	ChatId     typ.ChatId
	ReplyId    typ.MessageId
	Author     string
	Text       string
	CreatedAt  time.Time
	LastEditAt time.Time
}

type Chat struct {
	Id                 typ.ChatId
	Name               string
	Admin              typ.UserId
	AdminName          string
	MemberCount        int64
	UnreadMessageCount int64
}

type User struct {
	Id      typ.UserId
	Name    string
	Email   typ.Email
	PwdHash typ.PwdHash
	Joined  time.Time
	Token   typ.JWE
}
