package userservice

import (
	dto "server/data/DTOs"
	ent "server/data/entities"
	i "server/interfaces"
	cred "server/services/authService/credentials"
	typ "server/types"
)

func NewUserService(l i.Logger, u i.UserRepository, c i.ChatService) *UserService {
	return &UserService{
		lgr:   l,
		usrR:  u,
		chatS: c,
	}
}

type UserService struct {
	lgr   i.Logger
	usrR  i.UserRepository
	chatS i.ChatService
}

func (u *UserService) GetUser(userId typ.UserId) (*ent.User, error) {
	u.lgr.LogFunctionInfo()
	return u.usrR.GetUser(userId)
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

	contact, err := u.usrR.FindUser(a.Email)
	if err != nil {
		return nil, err
	}

	if contact == nil {
		return nil, nil
	}

	return u.usrR.AddContact(typ.ContactId(contact.Id), contact.Name, contact.Email, a.UserId)
}

func (u *UserService) GetContacts(userId typ.UserId) ([]ent.Contact, error) {
	u.lgr.LogFunctionInfo()
	return u.usrR.GetContacts(userId)
}

func (u *UserService) GetContact(chatId typ.ChatId, userId typ.UserId) (*ent.Contact, error) {
	u.lgr.LogFunctionInfo()
	return u.usrR.GetContact(chatId, userId)

}
