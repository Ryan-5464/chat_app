package chatservice

import (
	"fmt"
	ent "server/data/entities"
	i "server/interfaces"
	typ "server/types"
	"time"
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

func (c *ChatService) NewChat(chat ent.Chat) (ent.Chat, error) {
	c.lgr.LogFunctionInfo()
	chat, err := c.chatR.NewChat(chat)
	if err != nil {
		return ent.Chat{}, fmt.Errorf("failed to create new chat: %w", err)
	}
	return chat, nil
}

func (c *ChatService) GetChats(userId typ.UserId) ([]ent.Chat, error) {
	c.lgr.LogFunctionInfo()

	chats, err := c.chatR.GetChats(userId)
	if err != nil {
		return []ent.Chat{}, fmt.Errorf("failed to get chats: %w", err)
	}

	return chats, nil
}

func testChats() []ent.Chat {
	chat1 := ent.Chat{
		Id:                 11,
		Name:               "test1",
		AdminId:            3,
		AdminName:          "alf",
		MemberCount:        4,
		UnreadMessageCount: 14,
		CreatedAt:          time.Now(),
	}
	chat2 := ent.Chat{
		Id:                 12,
		Name:               "test2",
		AdminId:            4,
		AdminName:          "derek",
		MemberCount:        3,
		UnreadMessageCount: 2,
		CreatedAt:          time.Now(),
	}
	return []ent.Chat{chat1, chat2}
}
