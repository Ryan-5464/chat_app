package repository

import (
	"fmt"
	"log"
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

func (m *MessageRepository) NewMessage(msgE entities.Message) (entities.Message, error) {
	m.lgr.LogFunctionInfo()
	msgM := messageModelFromEntity(msgE)

	if m.dbS == nil {
		log.Fatal("database service is nil")
	}

	newMsg, err := m.dbS.NewMessage(msgM)
	if err != nil {
		return entities.Message{}, fmt.Errorf("failed to create new message: %w", err)
	}

	newMsgE := messageEntityFromModel(newMsg)

	return newMsgE, nil
}

func (m *MessageRepository) GetMessages(msgIds []typ.MessageId) {
	m.lgr.LogFunctionInfo()
	m.dbS.GetMessages(msgIds)
}

func (m *MessageRepository) GetChatMessages(chatId typ.ChatId) ([]entities.Message, error) {
	m.lgr.LogFunctionInfo()

	msgMs, err := m.dbS.GetChatMessages(chatId)
	if err != nil {
		return []entities.Message{}, fmt.Errorf("failed to get messages: %w", err)
	}

	msgEs := messageEntitiesFromModels(msgMs)
	return msgEs, nil
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
