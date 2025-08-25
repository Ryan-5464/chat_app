package userservice

import (
	ent "server/data/entities"
	i "server/interfaces"
	cred "server/services/auth/credentials"
	typ "server/types"
	"server/util"
)

func NewUserService(u i.UserRepository) *UserService {
	return &UserService{
		usrR: u,
	}
}

type UserService struct {
	usrR i.UserRepository
}

func (u *UserService) GetUser(userId typ.UserId) (*ent.User, error) {
	util.Log.FunctionInfo()
	return u.usrR.GetUser(userId)
}

func (u *UserService) GetUserByEmail(email cred.Email) (*ent.User, error) {
	util.Log.FunctionInfo()
	return u.usrR.GetUserByEmail(email)
}

func (u *UserService) GetUsers(userIds []typ.UserId) ([]ent.User, error) {
	util.Log.FunctionInfo()

	var users []ent.User

	users, err := u.usrR.GetUsers(userIds)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (u *UserService) GetChatUsers(chatId typ.ChatId) ([]ent.User, error) {
	util.Log.FunctionInfo()

	var users []ent.User

	users, err := u.usrR.GetChatUsers(chatId)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (u *UserService) NewUser(name string, email cred.Email, pwdHash cred.PwdHash) (*ent.User, error) {
	util.Log.FunctionInfo()

	user, err := u.usrR.NewUser(name, email, pwdHash)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) FindUsers(emails []cred.Email) ([]ent.User, error) {
	util.Log.FunctionInfo()

	var users []ent.User

	users, err := u.usrR.FindUsers(emails)
	if err != nil {
		return users, err
	}

	return users, nil
}

func (u *UserService) AddContact(email cred.Email, userId typ.UserId) (*ent.Contact, error) {
	util.Log.FunctionInfo()

	contact, err := u.usrR.FindUser(email)
	if err != nil {
		return nil, err
	}

	if contact == nil {
		return nil, nil
	}

	return u.usrR.AddContact(typ.ContactId(contact.Id), contact.Name, contact.Email, userId)
}

func (u *UserService) EditUserName(name string, userId typ.UserId) error {
	util.Log.FunctionInfo()
	return u.usrR.EditUserName(name, userId)
}

func (u *UserService) GetContacts(userId typ.UserId) ([]ent.Contact, error) {
	util.Log.FunctionInfo()
	return u.usrR.GetContacts(userId)
}

func (u *UserService) GetContact(chatId typ.ChatId, userId typ.UserId) (*ent.Contact, error) {
	util.Log.FunctionInfo()
	return u.usrR.GetContact(chatId, userId)

}

func (u *UserService) RemoveContact(contactId typ.ContactId, userId typ.UserId) error {
	util.Log.FunctionInfo()
	return u.usrR.RemoveContact(contactId, userId)
}
