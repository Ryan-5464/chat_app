package userservice

import (
	"errors"
	"fmt"
	dto "server/data/DTOs"
	ent "server/data/entities"
	i "server/interfaces"
	cred "server/services/authService/credentials"
	typ "server/types"
)

func NewUserService(l i.Logger, u i.UserRepository, c i.ChatRepository) *UserService {
	return &UserService{
		lgr:   l,
		usrR:  u,
		chatR: c,
	}
}

type UserService struct {
	lgr   i.Logger
	usrR  i.UserRepository
	chatR i.ChatRepository
}

func (u *UserService) GetUsers(userIds []typ.UserId) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()

	var users []ent.User

	users, err := u.usrR.GetUsers(userIds)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (u *UserService) GetChatUsers(chatId typ.ChatId) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()

	var users []ent.User

	users, err := u.usrR.GetChatUsers(chatId)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (u *UserService) NewUser(userReq dto.NewUserInput) (*ent.User, error) {
	u.lgr.LogFunctionInfo()

	user, err := u.usrR.NewUser(userReq.Name, userReq.Email, userReq.PwdHash)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) FindUsers(emails []cred.Email) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()

	var users []ent.User

	users, err := u.usrR.FindUsers(emails)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (u *UserService) AddContact(a dto.AddContactInput) (*ent.Contact, error) {
	u.lgr.LogFunctionInfo()

	user, err := u.usrR.FindUser(a.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, err
	}

	c := ent.Contact{
		Id:    user.Id,
		Name:  user.Name,
		Email: user.Email,
	}

	contact, err := u.usrR.AddContact(c, a.UserId)
	if err != nil {
		return nil, err
	}

	chatName := fmt.Sprintf("privateChat%v", a.UserId)
	chat, err := u.chatR.NewChat(chatName, a.UserId)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, errors.New("new chat missing")
	}

	contact.ContactChatId = chat.Id

	if err := u.chatR.NewMember(chat.Id, contact.Id); err != nil {
		return nil, err
	}

	return contact, nil
}

func (u *UserService) GetContacts(userId typ.UserId) ([]ent.Contact, error) {
	u.lgr.LogFunctionInfo()
	return u.usrR.GetContacts(userId)
}
