package chatservice

import (
	ent "server/data/entities"
	i "server/interfaces"
	cred "server/services/authService/credentials"
	typ "server/types"
)

func NewChatService(lgr i.Logger, c i.ChatRepository, u i.UserService) *ChatService {
	return &ChatService{
		lgr:   lgr,
		chatR: c,
		userS: u,
	}
}

type ChatService struct {
	lgr   i.Logger
	chatR i.ChatRepository
	userS i.UserService
}

func (c *ChatService) NewChat(chatName string, adminId typ.UserId) (*ent.Chat, error) {
	c.lgr.LogFunctionInfo()
	return c.chatR.NewChat(chatName, adminId)
}

func (c *ChatService) GetChats(userId typ.UserId) ([]ent.Chat, error) {
	c.lgr.LogFunctionInfo()
	return c.chatR.GetChats(userId)
}

func (c *ChatService) AddMember(email cred.Email, chatId typ.ChatId) (typ.UserId, error) {
	c.lgr.LogFunctionInfo()

	user, err := c.userS.GetUserByEmail(email)
	if err != nil {
		return 0, err
	}

	if err := c.chatR.NewMember(chatId, user.Id); err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (c *ChatService) GetChatMember(chatId typ.ChatId, userId typ.UserId) (*ent.Member, error) {
	c.lgr.LogFunctionInfo()

	memberships, err := c.chatR.GetChatMemberships(userId)
	if err != nil {
		return nil, err
	}

	var member ent.Member
	for _, m := range memberships {
		if m.ChatId == chatId {
			member = m
			break
		}
	}

	user, err := c.userS.GetUser(member.UserId)
	if err != nil {
		return nil, err
	}

	member.Name = user.Name
	member.Email = user.Email
	return &member, nil
}

func (c *ChatService) GetChatMembers(chatId typ.ChatId) ([]ent.Member, error) {
	c.lgr.LogFunctionInfo()

	members, err := c.chatR.GetMembers(chatId)
	if err != nil {
		return []ent.Member{}, err
	}

	memberIds := []typ.UserId{}
	for _, member := range members {
		memberIds = append(memberIds, member.UserId)
	}

	users, err := c.userS.GetUsers(memberIds)
	if err != nil {
		return []ent.Member{}, err
	}

	userIdMap := make(map[typ.UserId]ent.User)
	for _, user := range users {
		userIdMap[user.Id] = user
	}

	for i := range members {
		user := userIdMap[members[i].UserId]
		members[i].Name = user.Name
		members[i].Email = user.Email
	}

	return members, nil
}

func (c *ChatService) LeaveChat(chatId typ.ChatId, userId typ.UserId) ([]ent.Chat, error) {
	c.lgr.LogFunctionInfo()

	if err := c.chatR.RemoveChatMember(chatId, userId); err != nil {
		c.lgr.LogFunctionInfo()
		return []ent.Chat{}, err
	}

	chat, err := c.chatR.GetChat(chatId)
	if err != nil {
		return []ent.Chat{}, err
	}

	if chat == nil || chat.AdminId != userId {
		return c.chatR.GetChats(userId)
	}

	members, err := c.chatR.GetMembers(chatId)
	if err != nil {
		return []ent.Chat{}, err
	}

	if len(members) == 0 {
		if err := c.chatR.DeleteChat(chatId); err != nil {
			return []ent.Chat{}, err
		}
		return c.chatR.GetChats(userId)
	}

	if err := c.chatR.NewChatAdmin(chatId, members[0].UserId); err != nil {
		return []ent.Chat{}, err
	}

	return c.chatR.GetChats(userId)
}

func (c *ChatService) EditChatName(newName string, chatId typ.ChatId, userId typ.UserId) error {
	c.lgr.LogFunctionInfo()

	isAdmin, err := c.chatR.VerifyChatAdmin(chatId, userId)
	if err != nil {
		return err
	}

	if !isAdmin {
		return err
	}

	return c.chatR.EditChatName(newName, chatId)
}
