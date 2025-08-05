package repository

import (
	"errors"
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

func (c *ChatRepository) GetChatMessages(chatId typ.ChatId) ([]ent.Message, error) {
	c.lgr.DLog(fmt.Sprintf("chatid %v", chatId))

	messages, err := c.dbS.GetChatMessages(chatId)
	if err != nil {
		return []ent.Message{}, err
	}

	if len(messages) == 0 {
		return []ent.Message{}, nil
	}

	return messageEntitiesFromModels(messages), nil
}

func (c *ChatRepository) NewMember(chatId typ.ChatId, userId typ.UserId) error {
	c.lgr.LogFunctionInfo()

	return c.dbS.NewMember(chatId, userId)
}

func (c *ChatRepository) NewChat(chatName string, adminId typ.UserId) (*ent.Chat, error) {
	c.lgr.LogFunctionInfo()

	chat, err := c.dbS.NewChat(chatName, adminId)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, errors.New("new chat missing")
	}

	if err := c.dbS.NewMember(chat.Id, chat.AdminId); err != nil {
		return nil, err
	}

	return chatModelToEntity(chat), nil
}

func (c *ChatRepository) GetChats(userId typ.UserId) ([]ent.Chat, error) {
	c.lgr.LogFunctionInfo()

	chats, err := c.dbS.GetUserChats(userId)
	if err != nil {
		return []ent.Chat{}, err
	}

	if len(chats) == 0 {
		return []ent.Chat{}, nil
	}

	return chatModelsToEntities(chats), nil
}

func (c *ChatRepository) NewContactChat(adminId typ.UserId, contactId typ.UserId) error {
	c.lgr.LogFunctionInfo()

	_, err := c.dbS.NewContactChat(adminId, contactId)
	if err != nil {
		return err
	}

	return nil
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
