package chatservice

import (
	"errors"
	ent "server/data/entities"
	i "server/interfaces"
	cred "server/services/auth/credentials"
	typ "server/types"
	"server/util"
)

func NewChatService(c i.ChatRepository, m i.MessageService, u i.UserService) *ChatService {
	return &ChatService{
		chatR: c,
		msgS:  m,
		userS: u,
	}
}

type ChatService struct {
	chatR i.ChatRepository
	msgS  i.MessageService
	userS i.UserService
}

func (c *ChatService) NewChat(chatName string, adminId typ.UserId) (*ent.Chat, error) {
	util.Log.FunctionInfo()
	return c.chatR.NewChat(chatName, adminId)
}

func (c *ChatService) GetChats(userId typ.UserId) ([]ent.Chat, error) {
	util.Log.FunctionInfo()
	chats, err := c.chatR.GetChats(userId)
	if err != nil {
		return []ent.Chat{}, err
	}

	for i := range chats {
		chats[i].UnreadMessageCount, err = c.GetUnreadMessageCount(chats[i].Id, userId)
		if err != nil {
			util.Log.Infof("Unable to get unread message count for chat id %v => error: %v", chats[i].Id, err)
		}
	}

	return chats, nil
}

func (c *ChatService) AddMember(email cred.Email, chatId typ.ChatId) (typ.UserId, error) {
	util.Log.FunctionInfo()

	user, err := c.userS.GetUserByEmail(email)
	if err != nil {
		return 0, err
	}

	if err := c.chatR.NewMember(chatId, user.Id); err != nil {
		return 0, err
	}

	latestMsgId, err := c.msgS.GetLatestChatMessageId(chatId)
	if err != nil {
		return 0, err
	}

	if err := c.msgS.UpdateLastReadMsgId(latestMsgId, chatId, user.Id); err != nil {
		return 0, err
	}

	return user.Id, nil
}

func (c *ChatService) GetChatMember(chatId typ.ChatId, userId typ.UserId) (*ent.Member, error) {
	util.Log.FunctionInfo()

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
	util.Log.FunctionInfo()

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

func (c *ChatService) RemoveMember(chatId typ.ChatId, userId typ.UserId, adminId typ.UserId) error {
	util.Log.FunctionInfo()

	isAdmin, err := c.chatR.VerifyChatAdmin(chatId, adminId)
	if err != nil {
		return err
	}

	if !isAdmin {
		return errors.New("user not authorized to remove members")
	}

	return c.chatR.RemoveChatMember(chatId, userId)
}

func (c *ChatService) LeaveChat(chatId typ.ChatId, userId typ.UserId) ([]ent.Chat, error) {
	util.Log.FunctionInfo()

	if err := c.chatR.RemoveChatMember(chatId, userId); err != nil {
		util.Log.FunctionInfo()
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
	util.Log.FunctionInfo()

	isAdmin, err := c.chatR.VerifyChatAdmin(chatId, userId)
	if err != nil {
		return err
	}

	if !isAdmin {
		return err
	}

	return c.chatR.EditChatName(newName, chatId)
}

func (c *ChatService) GetUnreadMessageCount(chatId typ.ChatId, userId typ.UserId) (int64, error) {
	util.Log.FunctionInfo()

	member, err := c.chatR.GetMember(chatId, userId)
	if err != nil {
		return 0, err
	}

	util.Log.Dbugf("MEMBER BEFORE => %v", member)
	if member.LastReadMsgId == 0 {
		return 0, nil
	}
	util.Log.Dbugf("MEMBER AFTER => %v", member)

	return c.chatR.GetUnreadMessageCount(member.LastReadMsgId, chatId)
}
