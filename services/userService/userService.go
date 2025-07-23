package userservice

import (
	"server/data/entities"
	i "server/interfaces"
	typ "server/types"
)

func NewUserService(u i.UserRepository) *UserService {
	return &UserService{
		usrR: u,
	}
}

type UserService struct {
	usrR i.UserRepository
}

func (u *UserService) GetUsers(chatId typ.ChatId) ([]entities.User, error) {
	return testUsers(chatId), nil
}
