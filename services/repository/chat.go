package repository

import (
	"errors"
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

func (c *ChatRepository) NewChat(newChat ent.Chat) (*ent.Chat, error) {
	c.lgr.LogFunctionInfo()

	chatModel := chatEntityToModel(newChat)
	newChatModel, err := c.dbS.NewChat(chatModel)
	if err != nil {
		return nil, err
	}

	if newChatModel == nil {
		return nil, errors.New("new chat missing")
	}

	if err := c.dbS.NewMember(newChatModel.Id, newChatModel.AdminId); err != nil {
		return nil, err
	}

	return chatModelToEntity(newChatModel), nil
}

func (c *ChatRepository) GetChats(userId typ.UserId) ([]ent.Chat, error) {
	c.lgr.LogFunctionInfo()

	var chats []ent.Chat

	chatModels, err := c.dbS.GetUserChats(userId)
	if err != nil {
		return chats, err
	}

	if len(chatModels) == 0 {
		return chats, nil
	}

	return chatModelsToEntities(chatModels), nil
}

func chatEntitiesToModels(chats []ent.Chat) []model.Chat {
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

func chatModelToEntity(chat *model.Chat) *ent.Chat {
	return &ent.Chat{
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
