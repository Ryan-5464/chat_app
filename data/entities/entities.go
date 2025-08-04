package entities

import (
	cred "server/services/authService/credentials"
	tkn "server/services/authService/jwetoken"
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
	AdminId            typ.UserId
	AdminName          string
	MemberCount        int64
	UnreadMessageCount int64
	CreatedAt          time.Time
}

type User struct {
	Id      typ.UserId
	Name    string
	Email   cred.Email
	PwdHash cred.PwdHash
	Joined  time.Time
	Token   tkn.JWE
}

func (u User) IdIsZero() bool {
	return u.Id == 0
}

type Contact struct {
	Id           typ.UserId
	Name         string
	Email        cred.Email
	KnownSince   time.Time
	OnlineStatus bool
}
