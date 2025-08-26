package repository

import (
	"errors"
	ent "server/data/entities"
	i "server/interfaces"
	model "server/services/db/SQL/models"
	typ "server/types"
	"server/util"
)

func NewChatRepository(dbS i.DbService) *ChatRepository {
	return &ChatRepository{
		dbS: dbS,
	}
}

type ChatRepository struct {
	dbS i.DbService
}

func (c *ChatRepository) GetChatMessages(chatId typ.ChatId) ([]ent.Message, error) {
	util.Log.FunctionInfo()

	messages, err := c.dbS.GetChatMessages(chatId)
	if err != nil {
		return []ent.Message{}, err
	}

	return messageEntitiesFromModels(messages), nil
}

func (c *ChatRepository) NewMember(chatId typ.ChatId, userId typ.UserId) error {
	util.Log.FunctionInfo()
	return c.dbS.CreateMember(chatId, userId)
}

func (c *ChatRepository) NewChat(chatName string, adminId typ.UserId) (*ent.Chat, error) {
	util.Log.FunctionInfo()

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

func (c *ChatRepository) GetChat(chatId typ.ChatId) (*ent.Chat, error) {
	util.Log.FunctionInfo()

	chat, err := c.dbS.GetChat(chatId)
	if err != nil {
		return nil, err
	}

	return chatModelToEntity(chat), nil
}

func (c *ChatRepository) GetChats(userId typ.UserId) ([]ent.Chat, error) {
	util.Log.FunctionInfo()

	memberships, err := c.dbS.GetMemberships(userId)
	if err != nil {
		return []ent.Chat{}, err
	}

	chatIds := getMembershipIds(memberships)

	chats, err := c.dbS.GetChats(chatIds)
	if err != nil {
		return []ent.Chat{}, err
	}

	return chatModelsToEntities(chats), nil
}

func (c *ChatRepository) DeleteChat(chatId typ.ChatId) error {
	util.Log.FunctionInfo()
	return c.dbS.DeleteChat(chatId)
}

func (c *ChatRepository) GetUnreadMessageCount(lastReadMsgId typ.MessageId, chatId typ.ChatId) (int64, error) {
	util.Log.FunctionInfo()
	return c.dbS.GetUnreadMessageCount(lastReadMsgId, chatId)
}

func (c *ChatRepository) GetMember(chatId typ.ChatId, userId typ.UserId) (*ent.Member, error) {
	util.Log.FunctionInfo()
	member, err := c.dbS.GetMember(chatId, userId)
	if err != nil {
		return nil, err
	}

	return memberEntityFromModel(member), nil
}

func (c *ChatRepository) GetMembers(chatId typ.ChatId) ([]ent.Member, error) {
	util.Log.FunctionInfo()

	members, err := c.dbS.GetMembers(chatId)
	if err != nil {
		return []ent.Member{}, err
	}

	return memberModelsToEntities(members), nil
}

func (c *ChatRepository) GetChatMemberships(userId typ.UserId) ([]ent.Member, error) {
	util.Log.FunctionInfo()
	memberships, err := c.dbS.GetMemberships(userId)
	if err != nil {
		return []ent.Member{}, err
	}

	return memberModelsToEntities(memberships), nil
}

func (c *ChatRepository) RemoveChatMember(chatId typ.ChatId, userId typ.UserId) error {
	util.Log.FunctionInfo()
	return c.dbS.DeleteMember(chatId, userId)
}

func (c *ChatRepository) NewChatAdmin(chatId typ.ChatId, newAdminId typ.UserId) error {
	util.Log.FunctionInfo()
	return c.dbS.UpdateChatAdmin(chatId, newAdminId)
}

func (c *ChatRepository) VerifyChatAdmin(chatId typ.ChatId, userId typ.UserId) (bool, error) {
	util.Log.FunctionInfo()

	chat, err := c.dbS.GetChat(chatId)
	if err != nil {
		return false, err
	}

	if chat == nil {
		return false, errors.New("Failed to find chat")
	}

	return chat.AdminId == userId, nil
}

func (c *ChatRepository) EditChatName(newName string, chatId typ.ChatId) error {
	util.Log.FunctionInfo()
	return c.dbS.UpdateChatName(newName, chatId)
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

func getMembershipIds(members []model.Member) []typ.ChatId {
	var chatIds []typ.ChatId
	for _, member := range members {
		chatIds = append(chatIds, member.ChatId)
	}
	return chatIds
}

func memberModelsToEntities(models []model.Member) []ent.Member {
	var members []ent.Member
	for _, model := range models {
		ent := ent.Member{
			ChatId:        model.ChatId,
			UserId:        model.UserId,
			LastReadMsgId: model.LastReadMsgId,
			Joined:        model.Joined,
		}
		members = append(members, ent)
	}
	return members
}

func memberEntityFromModel(model *model.Member) *ent.Member {
	return &ent.Member{
		ChatId:        model.ChatId,
		UserId:        model.UserId,
		LastReadMsgId: model.LastReadMsgId,
		Joined:        model.Joined,
	}
}
