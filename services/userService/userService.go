package userservice

import (
	"fmt"
	ent "server/data/entities"
	i "server/interfaces"
	typ "server/types"
)

func NewUserService(l i.Logger, u i.UserRepository) *UserService {
	return &UserService{
		lgr:  l,
		usrR: u,
	}
}

type UserService struct {
	lgr  i.Logger
	usrR i.UserRepository
}

func (u *UserService) GetUser(uid typ.UserId) (ent.User, error) {
	u.lgr.LogFunctionInfo()

	user, err := u.usrR.GetUser(uid)
	if err != nil {
		return ent.User{}, fmt.Errorf("faield to get user: %w", err)
	}

	return user, nil
}

func (u *UserService) GetUsers(chatId typ.ChatId) ([]ent.User, error) {
	return nil, nil
}

func (u *UserService) GetUsersForChat(chatId typ.ChatId) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()
	usrEs, err := u.usrR.GetUsersForChat(chatId)
	if err != nil {
		return nil, err
	}

	return usrEs, nil
}

func (u *UserService) NewUser(usr ent.User) (ent.User, error) {
	u.lgr.LogFunctionInfo()
	usr, err := u.usrR.NewUser(usr)
	if err != nil {
		return ent.User{}, fmt.Errorf("failed to crate new user: %w", err)
	}
	return usr, nil
}

func (u *UserService) FindUser(usr ent.User) (ent.User, error) {
	u.lgr.LogFunctionInfo()
	usr, err := u.usrR.FindUser(usr)
	if err != nil {
		return ent.User{}, fmt.Errorf("failed to crate new user: %w", err)
	}
	return usr, nil
}

func (u *UserService) AddFriend(friend ent.Friend, userId typ.UserId) (ent.Friend, error) {
	u.lgr.LogFunctionInfo()

	user, err := u.usrR.FindUserByEmail(friend.Email)
	if err != nil {
		return ent.Friend{}, err
	}

	if user.IdIsZero() {
		return ent.Friend{}, err
	}

	friend, err = u.usrR.AddFriend(friend, userId)
	if err != nil {
		return ent.Friend{}, err
	}
	return friend, nil
}
