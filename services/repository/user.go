package repository

import (
	"fmt"
	"server/data/entities"
	i "server/interfaces"
	model "server/services/dbService/SQL/models"
	typ "server/types"
)

func NewUserRepository(dbS i.DbService) *UserRepository {
	return &UserRepository{dbS: dbS}
}

type UserRepository struct {
	dbS i.DbService
}

func (u *UserRepository) GetUsers(chatId typ.ChatId) ([]entities.User, error) {
	usrMs, err := u.dbS.GetUsers(chatId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user models from database service: %w", err)
	}

	return userEntitiesFromModels(usrMs), nil
}

func (u *UserRepository) NewUser(usrE entities.User) error {

	usrM := userModelFromEntity(usrE)
	if err := u.dbS.NewUser(usrM); err != nil {
		return fmt.Errorf("failed to create new user: %w", err)
	}

	return nil
}

func userEntitiesFromModels(usrMs []model.User) []entities.User {
	var usrEs []entities.User
	for _, usrM := range usrMs {
		usrE := userEntityFromModel(usrM)
		usrEs = append(usrEs, usrE)
	}
	return usrEs
}

func userEntityFromModel(m model.User) entities.User {
	return entities.User{
		Id:      m.Id,
		Name:    m.Name,
		Email:   m.Email,
		PwdHash: m.PwdHash,
		Joined:  m.Joined,
	}
}

func userModelFromEntity(u entities.User) model.User {
	return model.User{
		Id:      u.Id,
		Name:    u.Name,
		Email:   u.Email,
		PwdHash: u.PwdHash,
		Joined:  u.Joined,
	}
}
