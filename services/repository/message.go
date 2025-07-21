package repository

import (
	i "server/interfaces"
)

func NewMessageRepository(dbS i.DbService) *MessageRepository {
	return &MessageRepository{dbS: dbS}
}

type MessageRepository struct {
	dbS i.DbService
}

func (u *MessageRepository) GetChats() {
	u.dbS.GetMessages()
}
