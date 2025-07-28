package repository

import (
	"fmt"
	ent "server/data/entities"
	i "server/interfaces"
	model "server/services/dbService/SQL/models"
	typ "server/types"
)

func NewChatRepository(lgr i.Logger, dbS i.DbService) *ChatRepository {
	return &ChatRepository{
		lgr: lgr,
		dbS: dbS,
	}
}

type ChatRepository struct {
	lgr i.Logger
	dbS i.DbService
}

func (c *ChatRepository) NewChat(chat ent.Chat) (ent.Chat, error) {
	chatM := chatEntityToModel(chat)
	newChatM, err := c.dbS.NewChat(chatM)
	if err != nil {
		return ent.Chat{}, fmt.Errorf("failed to create new chat: %w", err)
	}

	if err := c.dbS.NewMember(newChatM.Id, newChatM.AdminId); err != nil {
		return ent.Chat{}, err
	}

	chatE := chatModelToEntity(newChatM)
	return chatE, nil
}

func (c *ChatRepository) GetChats(userId typ.UserId) ([]ent.Chat, error) {
	c.lgr.LogFunctionInfo()
	chats, err := c.dbS.GetChats(userId)
	if err != nil {
		return []ent.Chat{}, fmt.Errorf("faied to get chats: %w", err)
	}
	return chatModelsToEntities(chats), nil
}

func chatentitiesToModels(chats []ent.Chat) []model.Chat {
	chatMdls := []model.Chat{}
	for _, chat := range chats {
		mdl := model.Chat{
			Id:        chat.Id,
			Name:      chat.Name,
			AdminId:   chat.AdminId,
			CreatedAt: chat.CreatedAt,
		}
		chatMdls = append(chatMdls, mdl)
	}
	return chatMdls
}

func chatModelsToEntities(chats []model.Chat) []ent.Chat {
	chatEnts := []ent.Chat{}
	for _, chat := range chats {
		ent := ent.Chat{
			Id:        chat.Id,
			Name:      chat.Name,
			AdminId:   chat.AdminId,
			CreatedAt: chat.CreatedAt,
		}
		chatEnts = append(chatEnts, ent)
	}
	return chatEnts
}

func chatModelToEntity(chat model.Chat) ent.Chat {
	return ent.Chat{
		Id:        chat.Id,
		Name:      chat.Name,
		AdminId:   chat.AdminId,
		CreatedAt: chat.CreatedAt,
	}
}

func chatEntityToModel(chat ent.Chat) model.Chat {
	return model.Chat{
		Id:        chat.Id,
		Name:      chat.Name,
		AdminId:   chat.AdminId,
		CreatedAt: chat.CreatedAt,
	}
}
