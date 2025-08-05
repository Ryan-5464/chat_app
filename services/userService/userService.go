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

	newUser := ent.User{
		Email:   userReq.Email,
		PwdHash: userReq.PwdHash,
		Name:    userReq.Name,
	}

	user, err := u.usrR.NewUser(newUser)
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

	contactEmails := []cred.Email{a.Email}
	users, err := u.usrR.FindUsers(contactEmails)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, err
	}

	c := ent.Contact{
		Id:    users[0].Id,
		Name:  users[0].Name,
		Email: users[0].Email,
	}

	contact, err := u.usrR.AddContact(c, a.UserId)
	if err != nil {
		return nil, err
	}

	chatName := fmt.Sprintf("privateChat%v", a.UserId)
	chat, err := u.chatR.NewChat(chatName, a.UserId, typ.Private)
	if err != nil {
		return nil, err
	}

	if chat == nil {
		return nil, errors.New("new chat missing")
	}

	contact.PrivateChatId = chat.Id

	if err := u.chatR.NewMember(chat.Id, contact.Id); err != nil {
		return nil, err
	}

	return contact, nil
}

func (u *UserService) GetContacts(userId typ.UserId) ([]ent.Contact, error) {
	u.lgr.LogFunctionInfo()

	var contacts []ent.Contact

	contacts, err := u.usrR.GetContacts(userId)
	if err != nil {
		return contacts, err
	}

	return contacts, nil
}
