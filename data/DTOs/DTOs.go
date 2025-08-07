package DTO

import (
	"encoding/json"
	"server/data/entities"
	"server/lib"
	cred "server/services/authService/credentials"
	typ "server/types"
)

type SwitchChatRequest struct {
	ChatId string `json:"ChatId"`
}

func (s *SwitchChatRequest) GetChatId() (typ.ChatId, error) {
	cid, err := lib.ConvertStringToInt64(s.ChatId)
	if err != nil {
		return typ.ChatId(0), err
	}
	return typ.ChatId(cid), nil
}

type SwitchChatResponse struct {
	NewActiveChatId typ.ChatId
	Messages        []entities.Message `json:"Messages"`
}

type NewChatRequest struct {
	Name string `json:"Name"`
}

type NewChatResponse struct {
	Chats           []entities.Chat    `json:"Chats"`
	Messages        []entities.Message `json:"Messages"`
	NewActiveChatId typ.ChatId         `json:"NewActiveChatId"`
}

type NewMessageRequest struct {
	ChatId  string `json:"ChatId"`
	ReplyId string `json:"ReplyId"`
	MsgText string `json:"MsgText"`
}

type WebsocketPayload struct {
	Type string          `json:"Type"`
	Data json.RawMessage `json:"Data"`
}

func (w *WebsocketPayload) ParseNewMessageRequest() (NewMessageRequest, error) {
	newMessageRequest := NewMessageRequest{}
	if err := json.Unmarshal(w.Data, &newMessageRequest); err != nil {
		return NewMessageRequest{}, err
	}
	return newMessageRequest, nil
}

type ResponsePayload struct {
	Type     string             `json:"Type"`
	Chats    []entities.Chat    `json:"Chats"`
	Messages []entities.Message `json:"Messages"`
}

type RenderChatPayload struct {
	UserId   typ.UserId         `json:"UserId"`
	Chats    []entities.Chat    `json:"Chats"`
	Messages []entities.Message `json:"Messages"`
	Contacts []entities.Contact `json:"Contacts"`
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
