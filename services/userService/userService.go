package userservice

import (
	"fmt"
	"log"
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
	log.Println(1)

	u.lgr.LogFunctionInfo()
	log.Println(2)

	user, err := u.usrR.GetUser(uid)
	if err != nil {
		return ent.User{}, fmt.Errorf("faield to get user: %w", err)
	}
	log.Println(3)

	return user, nil
}

func (u *UserService) GetUsers(chatId typ.ChatId) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()
	return []ent.User{}, nil
	// return testUsers(chatId), nil
}

func (u *UserService) NewUser(usr ent.User) (ent.User, error) {
	u.lgr.LogFunctionInfo()
	usr, err := u.usrR.NewUser(usr)
	if err != nil {
		return ent.User{}, fmt.Errorf("failed to crate new user: %w", err)
	}
	return usr, nil
}
