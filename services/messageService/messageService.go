package messageservice

import (
	"server/data/entities"
	"time"
)

type MessageService struct {
}

func (m *MessageService) GetMessages() ([]entities.Message, error) {
	return testMessages(), nil
}

func testMessages() []entities.Message {
	message1 := entities.Message{
		Id:         1,
		UserId:     3,
		ChatId:     1,
		ReplyId:    0,
		Author:     "alf",
		Text:       "hello",
		CreatedAt:  time.Now(),
		LastEditAt: time.Now(),
	}
	message2 := entities.Message{
		Id:         2,
		UserId:     3,
		ChatId:     1,
		ReplyId:    0,
		Author:     "alf",
		Text:       "there",
		CreatedAt:  time.Now(),
		LastEditAt: time.Now(),
	}
	return []entities.Message{message1, message2}
}
