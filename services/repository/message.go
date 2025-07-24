package repository

import (
	"fmt"
	"server/data/entities"
	i "server/interfaces"
	model "server/services/dbService/SQL/models"
	typ "server/types"
)

func NewMessageRepository(dbS i.DbService) *MessageRepository {
	return &MessageRepository{dbS: dbS}
}

type MessageRepository struct {
	dbS i.DbService
}

func (u *MessageRepository) NewMessage(msgE entities.Message) (entities.Message, error) {
	msgM := messageModelFromEntity(msgE)

	newMsg, err := u.dbS.NewMessage(msgM)
	if err != nil {
		return entities.Message{}, fmt.Errorf("failed to create new message: %w", err)
	}

	newMsgE := messageEntityFromModel(newMsg)

	return newMsgE, nil
}

func (u *MessageRepository) GetMessages(msgIds []typ.MessageId) {
	u.dbS.GetMessages(msgIds)
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
