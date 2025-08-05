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

func (m *MessageRepository) NewMessage(newMsg entities.Message) (*entities.Message, error) {
	m.lgr.LogFunctionInfo()
	msgModel := messageModelFromEntity(newMsg)

	msg, err := m.dbS.NewMessage(msgModel)
	if err != nil {
		return nil, err
	}

	if msg == nil {
		return nil, errors.New("new message missing!")
	}

	msgEnt := messageEntityFromModel(*msg)

	return &msgEnt, nil
}

func (m *MessageRepository) GetMessages(msgIds []typ.MessageId) {
	m.lgr.LogFunctionInfo()
	m.dbS.GetMessages(msgIds)
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

func messageEntitiesFromModels(msgMs []model.Message) []entities.Message {
	var msgEs []entities.Message
	for _, msg := range msgMs {
		usrE := messageEntityFromModel(msg)
		msgEs = append(msgEs, usrE)
	}
	return msgEs
}

func messageEntityFromModel(m model.Message) entities.Message {
	return entities.Message{
		Id:         m.Id,
		UserId:     m.UserId,
		ChatId:     m.ChatId,
		ReplyId:    m.ReplyId,
		Text:       m.Text,
		CreatedAt:  m.CreatedAt,
		LastEditAt: m.LastEditAt,
	}
}

func messageModelFromEntity(m entities.Message) model.Message {
	return model.Message{
		Id:         m.Id,
		UserId:     m.UserId,
		ChatId:     m.ChatId,
		ReplyId:    m.ReplyId,
		Text:       m.Text,
		CreatedAt:  m.CreatedAt,
		LastEditAt: m.LastEditAt,
	}
}
