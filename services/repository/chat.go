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

func (c *ChatRepository) GetChatMessages(chatId typ.ChatId) ([]ent.Message, error) {
	c.lgr.DLog(fmt.Sprintf("chatid %v", chatId))

	messages, err := c.dbS.GetChatMessages(chatId)
	if err != nil {
		return []ent.Message{}, err
	}

	return messageEntitiesFromModels(messages), nil
}

func (c *ChatRepository) NewMember(chatId typ.ChatId, userId typ.UserId) error {
	c.lgr.LogFunctionInfo()
	return c.dbS.CreateMember(chatId, userId)
}

func (c *ChatRepository) NewChat(chatName string, adminId typ.UserId) (*ent.Chat, error) {
	c.lgr.LogFunctionInfo()

	lastInsertId, err := c.dbS.CreateChat(chatName, adminId)
	if err != nil {
		return nil, err
	}

	chat, err := c.dbS.GetChat(typ.ChatId(lastInsertId))
	if err != nil {
		return nil, err
	}

	if err := c.dbS.CreateMember(chat.Id, chat.AdminId); err != nil {
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

	return chatModelsToEntities(chats), nil
}

func chatModelsToEntities(chats []model.Chat) []ent.Chat {
	if len(chats) == 0 {
		return []ent.Chat{}
	}

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
	if chat == nil {
		return nil
	}

	return &ent.Chat{
		Id:        chat.Id,
		Name:      chat.Name,
		AdminId:   chat.AdminId,
		CreatedAt: chat.CreatedAt,
	}
}
