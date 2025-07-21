package chatservice

import (
	"server/data/entities"
	"time"
)

type ChatService struct {
}

func (m *ChatService) GetChats() ([]entities.Chat, error) {
	return testChats(), nil
}

func testChats() []entities.Chat {
	chat1 := entities.Chat{
		Id:                 1,
		Name:               "test1",
		AdminId:            3,
		AdminName:          "alf",
		MemberCount:        4,
		UnreadMessageCount: 14,
		CreatedAt:          time.Now(),
	}
	chat2 := entities.Chat{
		Id:                 2,
		Name:               "test2",
		AdminId:            4,
		AdminName:          "derek",
		MemberCount:        3,
		UnreadMessageCount: 2,
		CreatedAt:          time.Now(),
	}
	return []entities.Chat{chat1, chat2}
}
