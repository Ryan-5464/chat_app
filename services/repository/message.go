package repository

import (
	"errors"
	"server/data/entities"
	i "server/interfaces"
	model "server/services/dbService/SQL/models"
	typ "server/types"
)

func NewMessageRepository(lgr i.Logger, dbS i.DbService) *MessageRepository {
	return &MessageRepository{
		lgr: lgr,
		dbS: dbS,
	}
}

type MessageRepository struct {
	lgr i.Logger
	dbS i.DbService
}

func (m *MessageRepository) NewMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (*entities.Message, error) {
	m.lgr.LogFunctionInfo()

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

func (m *MessageRepository) NewContactMessage(userId typ.UserId, chatId typ.ChatId, replyId *typ.MessageId, text string) (*entities.Message, error) {
	m.lgr.LogFunctionInfo()

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

func (m *MessageRepository) GetChatMessages(chatId typ.ChatId) ([]entities.Message, error) {
	m.lgr.LogFunctionInfo()

	messages, err := m.dbS.GetChatMessages(chatId)
	if err != nil {
		return []entities.Message{}, err
	}

	if len(messages) == 0 {
		return []entities.Message{}, nil
	}

	return messageEntitiesFromModels(messages), nil
}

func messageEntitiesFromModels(messages []model.Message) []entities.Message {
	if len(messages) == 0 {
		return []entities.Message{}
	}

	var msgEs []entities.Message
	for _, m := range messages {
		usrE := entities.Message{
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

func messageEntityFromModel(m *model.Message) *entities.Message {
	if m == nil {
		return nil
	}

	return &entities.Message{
		Id:         m.Id,
		UserId:     m.UserId,
		ChatId:     m.ChatId,
		ReplyId:    m.ReplyId,
		Text:       m.Text,
		CreatedAt:  m.CreatedAt,
		LastEditAt: m.LastEditAt,
	}
}
