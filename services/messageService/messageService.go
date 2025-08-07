package messageservice

import (
	"fmt"
	"log"
	dto "server/data/DTOs"
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

func (m *MessageService) HandleNewContactMessage(mi dto.NewMessageInput) error {
	m.lgr.LogFunctionInfo()

	msg, err := m.msgR.NewContactMessage(mi.UserId, mi.ChatId, mi.ReplyId, mi.Text)
	if err != nil {
		return fmt.Errorf("failed to create new message: %w", err)
	}

	user, err := m.usrS.GetUser(msg.UserId)
	if err != nil {
		return err
	}

	msg.Author = user.Name

	contact, err := m.usrS.GetContact(msg.ChatId, mi.UserId)
	if err != nil {
		return err
	}

	usrConns := make(map[typ.UserId]i.Socket)
	usrConns[typ.UserId(contact.Id)] = m.connS.GetConnection(typ.UserId(contact.Id))
	usrConns[user.Id] = m.connS.GetConnection(user.Id)

	for userId, conn := range usrConns {
		if err := m.BroadcastMessage(userId, conn, *msg); err != nil {
			return err
		}
	}

	return nil
}

func (m *MessageService) HandleNewMessage(mi dto.NewMessageInput) error {
	m.lgr.LogFunctionInfo()

	msg, err := m.msgR.NewMessage(mi.UserId, mi.ChatId, mi.ReplyId, mi.Text)
	if err != nil {
		return fmt.Errorf("failed to create new message: %w", err)
	}

	user, err := m.usrS.GetUser(msg.UserId)
	if err != nil {
		return err
	}

	msg.Author = user.Name

	users, err := m.usrS.GetChatUsers(msg.ChatId)
	if err != nil {
		return err
	}

	usrConns := make(map[typ.UserId]i.Socket)
	for _, u := range users {
		conn := m.connS.GetConnection(u.Id)
		usrConns[u.Id] = conn
	}
	usrConns[user.Id] = m.connS.GetConnection(user.Id)

	for userId, conn := range usrConns {
		if err := m.BroadcastMessage(userId, conn, *msg); err != nil {
			return err
		}
		m.lgr.DLog("->>>> RESPONSE SENT")
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
	return m.msgR.GetChatMessages(chatId)
}
