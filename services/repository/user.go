package repository

import (
	i "server/interfaces"
)

func NewUserRepository(dbS i.DbService) *UserRepository {
	return &UserRepository{dbS: dbS}
}

type UserRepository struct {
	dbS i.DbService
}

func (u *UserRepository) GetUsers() {
	u.dbS.GetUsers()
}
