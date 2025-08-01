package DTO

import (
	"encoding/json"
	"server/data/entities"
	"server/lib"
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
	Chats           []entities.Chat
	Messages        []entities.Message
	NewActiveChatId typ.ChatId
}

type NewMessageRequest struct {
	ChatId  string `json:"ChatId"`
	ReplyId string `json:"ReplyId"`
	MsgText string `json:"MsgText"`
}

func (n *NewMessageRequest) ToMessageEntity(userId typ.UserId) (entities.Message, error) {
	chatId, err := lib.ConvertStringToInt64(n.ChatId)
	if err != nil {
		return entities.Message{}, err
	}

	var replyId int64
	if n.ReplyId != "" {
		replyId, err = lib.ConvertStringToInt64(n.ReplyId)
		if err != nil {
			return entities.Message{}, err
		}
	}

	msgE := entities.Message{
		UserId:  userId,
		ChatId:  typ.ChatId(chatId),
		ReplyId: typ.MessageId(replyId),
		Text:    n.MsgText,
	}

	return msgE, nil

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
	Chats    []entities.Chat    `json:"Chats"`
	Messages []entities.Message `json:"Messages"`
}

type ErrorResponse struct {
	NoError      bool   `json:"NoError"`
	ErrorMessage string `json:"ErrorMessage"`
}

type LoginRequest struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}
