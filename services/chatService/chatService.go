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

func (c *ChatService) NewChat(chatName string, adminId typ.UserId) (*ent.Chat, error) {
	c.lgr.LogFunctionInfo()
	return c.chatR.NewChat(chatName, adminId)
}

func (c *ChatService) GetChats(userId typ.UserId) ([]ent.Chat, error) {
	c.lgr.LogFunctionInfo()
	return c.chatR.GetChats(userId)
}
