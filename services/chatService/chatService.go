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
	return c.chatR.NewChat(newChat)
}

func (c *ChatService) GetChats(userId typ.UserId) ([]ent.Chat, error) {
	c.lgr.LogFunctionInfo()

	var chats []ent.Chat

	chats, err := c.chatR.GetChats(userId)
	if err != nil {
		return chats, err
	}

	if len(chats) == 0 {
		return chats, nil
	}

	return chats, nil
}
