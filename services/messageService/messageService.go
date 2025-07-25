package messageservice

import (
	"encoding/json"
	"fmt"
	"log"
	"server/data/entities"
	i "server/interfaces"
	typ "server/types"
	"time"
)

func NewMessageService(lgr i.Logger, m i.MessageRepository, u i.UserService, c i.ConnectionService) *MessageService {
	return &MessageService{
		lgr:   lgr,
		msgR:  m,
		usrS:  u,
		connS: c,
	}
}

type MessageService struct {
	lgr   i.Logger
	msgR  i.MessageRepository
	usrS  i.UserService
	connS i.ConnectionService
}

func (m *MessageService) HandleNewMessage(msg entities.Message) error {
	m.lgr.LogFunctionInfo()
	log.Println(1)

	msg, err := m.msgR.NewMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to create new message: %w", err)
	}
	log.Println(2)

	usr, err := m.usrS.GetUser(msg.UserId)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	log.Println(3)

	msg.Author = usr.Name
	log.Println(4)

	usrs, err := m.usrS.GetUsers(msg.ChatId)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}
	log.Println(5)

	usrConns := make(map[typ.UserId]i.Socket)
	for _, usr := range usrs {
		conn := m.connS.GetConnection(usr.Id)
		usrConns[usr.Id] = conn
	}
	log.Println(6)

	for userId, conn := range usrConns {
		if err := m.BroadcastMessage(userId, conn, msg); err != nil {
			return fmt.Errorf("failed to braodcast message: %w", err)
		}
	}
	log.Println(7)

	return nil
}

func (m *MessageService) BroadcastMessage(userId typ.UserId, conn i.Socket, msg entities.Message) error {
	m.lgr.LogFunctionInfo()
	messages := []entities.Message{msg}

	payload := struct {
		Chats    []entities.Chat
		Messages []entities.Message
	}{
		Chats:    nil,
		Messages: messages,
	}

	msgPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Failed to serialize data: %w", err)
	}

	err = conn.WriteJSON(msgPayload)
	if err != nil {
		return fmt.Errorf("failed to write to websockt connection: %w", err)
	}

	return nil
}

func (m *MessageService) GetMessages(chatId typ.ChatId) ([]entities.Message, error) {
	m.lgr.LogFunctionInfo()
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
