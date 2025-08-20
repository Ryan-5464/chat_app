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
		if conn == nil {
			m.lgr.Log(fmt.Sprintf("user is offline for userId %v", userId))
			continue
		}
		if err := m.BroadcastMessage(userId, conn, *msg); err != nil {
			m.lgr.LogError(fmt.Errorf(":: failed to broadcast message %v", err))
			return err
		}
	}
	log.Println(":: broadcast successful")

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
		if conn == nil {
			m.lgr.Log(fmt.Sprintf("connection is nil for userId %v!", userId))
			continue
		}
		if err := m.BroadcastMessage(userId, conn, *msg); err != nil {
			return err
		}
		m.lgr.DLog("->>>> RESPONSE SENT")
	}

	return nil
}

func (m *MessageService) BroadcastMessage(userId typ.UserId, conn i.Socket, msg entities.Message) error {
	m.lgr.LogFunctionInfo()
	log.Println(":: broadcast message ", msg)

	msg.IsUserMessage = msg.UserId == userId

	messages := []entities.Message{msg}

	payload := struct {
		Type     int
		Chats    []entities.Chat
		Messages []entities.Message
	}{
		Type:     1,
		Chats:    nil,
		Messages: messages,
	}

	log.Println(":: payload ", payload)

	err := conn.WriteJSON(payload)
	if err != nil {
		m.lgr.LogError(fmt.Errorf("failed to write to websocket connection: %w", err))
		return err
	}
	log.Println(":: json writing sucessful")

	return nil
}

func (m *MessageService) GetChatMessages(chatId typ.ChatId, userId typ.UserId) ([]entities.Message, error) {
	m.lgr.LogFunctionInfo()
	messages, err := m.msgR.GetChatMessages(chatId)
	if err != nil {
		return []entities.Message{}, err
	}

	userIds := []typ.UserId{}
	for _, message := range messages {
		userIds = append(userIds, message.UserId)
	}

	uniqueUserIds := getUniqueUserIdsFromMessages(userIds)

	users, err := m.usrS.GetUsers(uniqueUserIds)
	if err != nil {
		return []entities.Message{}, err
	}

	m.lgr.DLog(fmt.Sprintf("USERS => %v", users))

	authorMap := make(map[typ.UserId]string)
	for _, user := range users {
		authorMap[user.Id] = user.Name
	}

	for i := range messages {
		messages[i].Author = authorMap[messages[i].UserId]
		messages[i].IsUserMessage = messages[i].UserId == userId
	}

	m.lgr.DLog(fmt.Sprintf("MESSAGES => %v", messages))

	return messages, nil

}

func (m *MessageService) GetContactMessages(chatId typ.ChatId, userId typ.UserId) ([]entities.Message, error) {
	m.lgr.LogFunctionInfo()
	messages, err := m.msgR.GetContactMessages(chatId)
	if err != nil {
		return []entities.Message{}, err
	}

	for i := range messages {
		messages[i].IsUserMessage = messages[i].UserId == userId
	}

	return messages, nil
}

func (m *MessageService) DeleteMessage(messageId typ.MessageId) error {
	m.lgr.LogFunctionInfo()
	return m.msgR.DeleteMessage(messageId)
}

func getUniqueUserIdsFromMessages(slice []typ.UserId) []typ.UserId {
	seen := make(map[typ.UserId]struct{})
	var result []typ.UserId

	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}
