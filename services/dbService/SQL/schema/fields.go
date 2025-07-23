package schema

type Field string

const (
	UserId        Field = "UserId"
	Name          Field = "Name"
	Email         Field = "Email"
	PwdHash       Field = "PwdHash"
	CreatedAt     Field = "CreatedAt"
	LastReadMsgId Field = "LastReadMsgId"
	ChatId        Field = "ChatId"
	AdminId       Field = "AdminId"
	LastMsgAt     Field = "LastMsgAt"
	MessageId     Field = "MessageId"
	Author        Field = "Author"
	MsgText       Field = "MsgText"
	ReplyId       Field = "ReplyId"
	LastEditAt    Field = "LastEditAt"
)
