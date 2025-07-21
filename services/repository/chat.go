package repository

import (
	i "server/interfaces"
)

func NewChatRepository(dbS i.DbService) *ChatRepository {
	return &ChatRepository{dbS: dbS}
}

type ChatRepository struct {
	dbS i.DbService
}

func (u *ChatRepository) GetChats() {
	u.dbS.GetChats()
}
