package repository

import (
	"errors"
	ent "server/data/entities"
	i "server/interfaces"
	model "server/services/db/SQL/models"
	typ "server/types"
	"server/util"
)

func NewMessageRepository(dbS i.DbService) *MessageRepository {
	return &MessageRepository{
		dbS: dbS,
	}
}

type MessageRepository struct {
	dbS i.DbService
}

func (m *MessageRepository) NewMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (*ent.Message, error) {
	util.Log.FunctionInfo()

	lastInsertId, err := m.dbS.CreateMessage(userId, chatId, replyId, text)
	if err != nil {
		return nil, err
	}

	messageModel, err := m.dbS.GetMessage(typ.MessageId(lastInsertId))

	if messageModel == nil {
		return nil, errors.New("new message missing")
	}

	return messageEntityFromModel(messageModel), nil
}

func (m *MessageRepository) NewContactMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (*ent.Message, error) {
	util.Log.FunctionInfo()

	lastInsertId, err := m.dbS.CreateContactMessage(userId, chatId, replyId, text)
	if err != nil {
		return nil, err
	}

	messageModel, err := m.dbS.GetContactMessage(typ.MessageId(lastInsertId))

	if messageModel == nil {
		return nil, errors.New("new message missing")
	}

	return messageEntityFromModel(messageModel), nil
}

func (m *MessageRepository) GetContactMessages(chatId typ.ChatId) ([]ent.Message, error) {
	util.Log.FunctionInfo()

	messages, err := m.dbS.GetContactMessages(chatId)
	if err != nil {
		return []ent.Message{}, err
	}

	if len(messages) == 0 {
		return []ent.Message{}, nil
	}

	return messageEntitiesFromModels(messages), nil
}

func (m *MessageRepository) GetChatMessages(chatId typ.ChatId) ([]ent.Message, error) {
	util.Log.FunctionInfo()

	messages, err := m.dbS.GetChatMessages(chatId)
	if err != nil {
		return []ent.Message{}, err
	}

	if len(messages) == 0 {
		return []ent.Message{}, nil
	}

	return messageEntitiesFromModels(messages), nil
}

func (m *MessageRepository) DeleteMessage(messageId typ.MessageId) error {
	util.Log.FunctionInfo()
	return m.dbS.DeleteMessage(messageId)
}

func (m *MessageRepository) EditMessage(msgText string, msgId typ.MessageId) error {
	util.Log.FunctionInfo()
	return m.dbS.UpdateMessage(msgText, msgId)
}

func (m *MessageRepository) GetMessage(msgId typ.MessageId) (*ent.Message, error) {
	util.Log.FunctionInfo()

	msgModel, err := m.dbS.GetMessage(msgId)
	if err != nil {
		return nil, err
	}

	if msgModel == nil {
		return nil, errors.New("no message found")
	}

	return messageEntityFromModel(msgModel), nil
}

func (m *MessageRepository) GetLatestChatMessageId(chatId typ.ChatId) (typ.MessageId, error) {
	util.Log.FunctionInfo()
	return m.dbS.GetLatestChatMessageId(chatId)
}

func (m *MessageRepository) UpdateLastReadMsgId(lastReadMsgId typ.MessageId, chatId typ.ChatId, userId typ.UserId) error {
	util.Log.FunctionInfo()
	return m.dbS.UpdateLastReadMsgId(lastReadMsgId, chatId, userId)
}

func (m *MessageRepository) GetLatestMessageId() (typ.MessageId, error) {
	util.Log.FunctionInfo()
	return m.dbS.GetLatestMessageId()
}

func messageEntitiesFromModels(messages []model.Message) []ent.Message {
	if len(messages) == 0 {
		return []ent.Message{}
	}

	var msgEs []ent.Message
	for _, m := range messages {
		usrE := ent.Message{
			Id:         m.Id,
			UserId:     m.UserId,
			ChatId:     m.ChatId,
			ReplyId:    m.ReplyId,
			Text:       m.Text,
			CreatedAt:  m.CreatedAt,
			LastEditAt: m.LastEditAt,
		}
		msgEs = append(msgEs, usrE)
	}
	return msgEs
}

func messageEntityFromModel(m *model.Message) *ent.Message {
	if m == nil {
		return nil
	}

	return &ent.Message{
		Id:         m.Id,
		UserId:     m.UserId,
		ChatId:     m.ChatId,
		ReplyId:    m.ReplyId,
		Text:       m.Text,
		CreatedAt:  m.CreatedAt,
		LastEditAt: m.LastEditAt,
	}
}
