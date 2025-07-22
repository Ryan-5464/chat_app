package messageservice

import (
	"log"
	"server/data/entities"
	typ "server/types"
	"time"
)

func NewMessageService() *MessageService {
	return &MessageService{}
}

type MessageService struct {
}

func (m *MessageService) GetMessages(chatId typ.ChatId) ([]entities.Message, error) {
	return testMessages(chatId), nil
}

func testMessages(chatId typ.ChatId) []entities.Message {
	var messages []entities.Message
	log.Println("chatId", chatId)
	switch int(chatId) {
	case 1:
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
		messages = append(messages, message1, message2)

	case 2:
		message1 := entities.Message{
			Id:         1,
			UserId:     3,
			ChatId:     2,
			ReplyId:    0,
			Author:     "alf",
			Text:       "chat",
			CreatedAt:  time.Now(),
			LastEditAt: time.Now(),
		}
		message2 := entities.Message{
			Id:         2,
			UserId:     3,
			ChatId:     2,
			ReplyId:    0,
			Author:     "alf",
			Text:       "changed",
			CreatedAt:  time.Now(),
			LastEditAt: time.Now(),
		}
		messages = append(messages, message1, message2)
	}
	return messages
}
