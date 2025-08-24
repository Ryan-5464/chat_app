package dto

import (
	"server/data/entities"
	"server/lib"
	cred "server/services/authService/credentials"
	typ "server/types"
)

type SwitchContactChatRequest struct {
	ContactChatId string `json:"ContactChatId"`
}

type SwitchContactChatResponse struct {
	ActiveContactChatId typ.ChatId         `json:"ActiveContactChatId"`
	Messages            []entities.Message `json:"Messages"`
}

type NewMessageRequest struct {
	ChatId  string `json:"ChatId"`
	ReplyId string `json:"ReplyId"`
	MsgText string `json:"MsgText"`
}

type ResponsePayload struct {
	Type     string             `json:"Type"`
	Chats    []entities.Chat    `json:"Chats"`
	Messages []entities.Message `json:"Messages"`
}

type RenderChatPayload struct {
	UserId       typ.UserId         `json:"UserId"`
	Chats        []entities.Chat    `json:"Chats"`
	Messages     []entities.Message `json:"Messages"`
	Contacts     []entities.Contact `json:"Contacts"`
	ActiveChatId typ.ChatId         `json:"ActiveChatId"`
}

type ErrorResponse struct {
	NoError      bool   `json:"NoError"`
	ErrorMessage string `json:"ErrorMessage"`
}

type LoginRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type AddContactRequest struct {
	Email string `json:"Email"`
}

type AddContactResponse struct {
	Contacts []entities.Contact `json:"Contacts"`
}

type NewUserInput struct {
	Email   cred.Email
	PwdHash cred.PwdHash
	Name    string
}

type AddContactInput struct {
	Email  cred.Email
	UserId typ.UserId
}

type NewMessageInput struct {
	UserId  typ.UserId
	ChatId  typ.ChatId
	Text    string
	ReplyId *typ.MessageId
}

type EditUserNameRequest struct {
	Name string `json:"Name"`
}

type EditUserNameResponse struct {
	Name string
}

type EditMessageRequest struct {
	MsgText   string `json:"MsgText"`
	MessageId string `json:"MessageId"`
	UserId    string `json:"UserId"`
}

type EditMessageResponse struct {
	MsgText string
}

type AddMemberToChatRequest struct {
	Email  string `json:"Email"`
	ChatId string `json:"ChatId"`
}

type AddMemberToChatResponse struct {
	Members []entities.Member
}
