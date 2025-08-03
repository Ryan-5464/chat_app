package messageservice

import (
	"fmt"
	"log"
	"server/data/entities"
	i "server/interfaces"
	typ "server/types"
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

func (m *MessageService) NewMessage(newMsg entities.Message) (*entities.Message, error) {
	m.lgr.LogFunctionInfo()
	return m.msgR.NewMessage(newMsg)
}

func (m *MessageService) HandleNewMessage(newMsg entities.Message) error {
	m.lgr.LogFunctionInfo()

	msg, err := m.msgR.NewMessage(newMsg)
	if err != nil {
		return fmt.Errorf("failed to create new message: %w", err)
	}

	userIds := []typ.UserId{msg.UserId}
	users, err := m.usrS.GetUsers(userIds)
	if err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	user := users[0]

	msg.Author = user.Name

	users, err = m.usrS.GetChatUsers(msg.ChatId)
	if err != nil {
		return err
	}

	usrConns := make(map[typ.UserId]i.Socket)
	for _, usr := range users {
		conn := m.connS.GetConnection(usr.Id)
		usrConns[usr.Id] = conn
	}

	for userId, conn := range usrConns {
		if err := m.BroadcastMessage(userId, conn, msg); err != nil {
			return fmt.Errorf("failed to braodcast message: %w", err)
		}
	}

	return nil
}

func (m *MessageService) BroadcastMessage(userId typ.UserId, conn i.Socket, msg entities.Message) error {
	m.lgr.LogFunctionInfo()
	messages := []entities.Message{msg}

	log.Println("BROADCASTMSG ", messages)

	payload := struct {
		Type     int
		Chats    []entities.Chat
		Messages []entities.Message
	}{
		Type:     1,
		Chats:    nil,
		Messages: messages,
	}

	err := conn.WriteJSON(payload)
	if err != nil {
		return fmt.Errorf("failed to write to websockt connection: %w", err)
	}

	return nil
}

func (m *MessageService) GetChatMessages(chatId typ.ChatId) ([]entities.Message, error) {
	m.lgr.LogFunctionInfo()

	messages, err := m.msgR.GetChatMessages(chatId)
	if err != nil {
		return []entities.Message{}, fmt.Errorf("failed to get messsages: %w", err)
	}

	return messages, nil
}
