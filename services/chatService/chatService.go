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

func (c *ChatService) LeaveChat(chatId typ.ChatId, userId typ.UserId) ([]ent.Chat, error) {
	c.lgr.LogFunctionInfo()

	if err := c.chatR.RemoveChatMember(chatId, userId); err != nil {
		c.lgr.LogFunctionInfo()
		return []ent.Chat{}, err
	}

	chat, err := c.chatR.GetChat(chatId)
	if err != nil {
		return []ent.Chat{}, err
	}

	if chat.AdminId != userId {
		return c.chatR.GetChats(userId)
	}

	members, err := c.chatR.GetMembers(chatId)
	if err != nil {
		return []ent.Chat{}, err
	}

	if len(members) == 0 {
		if err := c.chatR.DeleteChat(chatId); err != nil {
			return []ent.Chat{}, err
		}
	}

	if err := c.chatR.NewChatAdmin(chatId, members[0].UserId); err != nil {
		return []ent.Chat{}, err
	}

	return c.chatR.GetChats(userId)
}
