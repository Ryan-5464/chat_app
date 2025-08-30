package messageservice

import (
	"fmt"
	ent "server/data/entities"
	i "server/interfaces"
	typ "server/types"
	"server/util"
)

func NewMessageService(m i.MessageRepository, u i.UserService, cn i.ConnectionService, c i.ChatService) *MessageService {
	return &MessageService{
		msgR:  m,
		usrS:  u,
		connS: cn,
		chatS: c,
	}
}

type MessageService struct {
	msgR  i.MessageRepository
	usrS  i.UserService
	connS i.ConnectionService
	chatS i.ChatService
}

// to avoid import cycle with chat service.
func (m *MessageService) SetChatService(c i.ChatService) {
	m.chatS = c
}

func (m *MessageService) HandleNewContactMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, msgTxt string) error {
	util.Log.FunctionInfo()

	msg, err := m.msgR.NewContactMessage(userId, chatId, replyId, msgTxt)
	if err != nil {
		return fmt.Errorf("failed to create new message: %w", err)
	}

	user, err := m.usrS.GetUser(msg.UserId)
	if err != nil {
		return err
	}

	msg.Author = user.Name

	contact, err := m.usrS.GetContact(msg.ChatId, userId)
	if err != nil {
		return err
	}

	usrConns := make(map[typ.UserId]i.Socket)
	usrConns[typ.UserId(contact.Id)] = m.connS.GetConnection(typ.UserId(contact.Id))
	usrConns[user.Id] = m.connS.GetConnection(user.Id)

	for userId, conn := range usrConns {
		if conn == nil {
			util.Log.Infof("user is offline for userId %v", userId)
			continue
		}
		if err := m.broadcastMessage(userId, nil, conn, *msg); err != nil {
			util.Log.Errorf(":: failed to broadcast message %v", err)
			return err
		}
	}
	util.Log.Dbug(":: broadcast successful")

	return nil
}

func (m *MessageService) HandleNewMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, msgTxt string) error {
	util.Log.FunctionInfo()

	msg, err := m.msgR.NewMessage(userId, chatId, replyId, msgTxt)
	if err != nil {
		return fmt.Errorf("failed to create new message: %w", err)
	}

	if err := m.msgR.UpdateLastReadMsgId(msg.Id, msg.ChatId, msg.UserId); err != nil {
		return err
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
			util.Log.Infof("connection is nil for userId %v!", userId)
			continue
		}

		chats, err := m.chatS.GetChats(userId)
		if err != nil {
			return err
		}

		if err := m.broadcastMessage(userId, chats, conn, *msg); err != nil {
			return err
		}
		util.Log.Dbug("->>>> RESPONSE SENT")
	}

	return nil
}

func (m *MessageService) broadcastMessage(userId typ.UserId, chats []ent.Chat, conn i.Socket, msg ent.Message) error {
	util.Log.FunctionInfo()

	msg.IsUserMessage = msg.UserId == userId

	messages := []ent.Message{msg}

	payload := struct {
		Type     int
		Chats    []ent.Chat
		Messages []ent.Message
	}{
		Type:     1,
		Chats:    chats,
		Messages: messages,
	}

	if err := conn.WriteJSON(payload); err != nil {
		util.Log.Errorf("failed to write to websocket connection: %v", err)
		return err
	}

	return nil
}

func (m *MessageService) GetLatestChatMessageId(chatId typ.ChatId) (typ.MessageId, error) {
	util.Log.FunctionInfo()
	return m.msgR.GetLatestChatMessageId(chatId)
}

func (m *MessageService) GetChatMessages(chatId typ.ChatId, userId typ.UserId) ([]ent.Message, error) {
	util.Log.FunctionInfo()
	messages, err := m.msgR.GetChatMessages(chatId)
	if err != nil {
		return []ent.Message{}, err
	}

	userIds := []typ.UserId{}
	for _, message := range messages {
		userIds = append(userIds, message.UserId)
	}

	uniqueUserIds := getUniqueUserIdsFromMessages(userIds)

	users, err := m.usrS.GetUsers(uniqueUserIds)
	if err != nil {
		return []ent.Message{}, err
	}

	authorMap := make(map[typ.UserId]string)
	for _, user := range users {
		authorMap[user.Id] = user.Name
	}

	for i := range messages {
		messages[i].Author = authorMap[messages[i].UserId]
		messages[i].IsUserMessage = messages[i].UserId == userId
	}

	return messages, nil

}

func (m *MessageService) GetContactMessages(chatId typ.ChatId, userId typ.UserId) ([]ent.Message, error) {
	util.Log.FunctionInfo()
	messages, err := m.msgR.GetContactMessages(chatId)
	if err != nil {
		return []ent.Message{}, err
	}

	userIds := []typ.UserId{}
	for _, message := range messages {
		userIds = append(userIds, message.UserId)
	}

	uniqueUserIds := getUniqueUserIdsFromMessages(userIds)

	users, err := m.usrS.GetUsers(uniqueUserIds)
	if err != nil {
		return []ent.Message{}, err
	}

	authorMap := make(map[typ.UserId]string)
	for _, user := range users {
		authorMap[user.Id] = user.Name
	}

	for i := range messages {
		messages[i].Author = authorMap[messages[i].UserId]
		messages[i].IsUserMessage = messages[i].UserId == userId
	}

	return messages, nil
}

func (m *MessageService) DeleteMessage(messageId typ.MessageId) error {
	util.Log.FunctionInfo()
	return m.msgR.DeleteMessage(messageId)
}

func (m *MessageService) EditMessage(msgText string, msgId typ.MessageId) (*ent.Message, error) {
	util.Log.FunctionInfo()

	if err := m.msgR.EditMessage(msgText, msgId); err != nil {
		return nil, err
	}

	return m.msgR.GetMessage(msgId)
}

func (m *MessageService) UpdateLastReadMsgId(lastReadMsgId typ.MessageId, chatId typ.ChatId, userId typ.UserId) error {
	util.Log.FunctionInfo()
	return m.msgR.UpdateLastReadMsgId(lastReadMsgId, chatId, userId)
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
