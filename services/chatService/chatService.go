package chatservice

import (
	ent "server/data/entities"
	i "server/interfaces"
	typ "server/types"
)

func NewChatService(lgr i.Logger, c i.ChatRepository) *ChatService {
	return &ChatService{
		lgr:   lgr,
		chatR: c,
	}
}

type ChatService struct {
	lgr   i.Logger
	chatR i.ChatRepository
}

func (c *ChatService) NewChat(newChat ent.Chat) (*ent.Chat, error) {
	c.lgr.LogFunctionInfo()
	return c.chatR.NewChat(newChat.Name, newChat.AdminId)
}

func (c *ChatService) GetChats(userId typ.UserId) ([]ent.Chat, error) {
	c.lgr.LogFunctionInfo()
	return c.chatR.GetChats(userId)
}
