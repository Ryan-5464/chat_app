package repository

import (
	"fmt"
	ent "server/data/entities"
	i "server/interfaces"
	model "server/services/dbService/SQL/models"
	typ "server/types"
)

func NewUserRepository(lgr i.Logger, dbS i.DbService) *UserRepository {
	return &UserRepository{
		lgr: lgr,
		dbS: dbS,
	}
}

type UserRepository struct {
	lgr i.Logger
	dbS i.DbService
}

func (u *UserRepository) FindUser(usr ent.User) (ent.User, error) {
	u.lgr.LogFunctionInfo()
	usrM, err := u.dbS.FindUser(usr.Email)
	if err != nil {
		return ent.User{}, fmt.Errorf("failed to find user in database service: %w", err)
	}
	return userEntityFromModel(usrM), nil
}

func (u *UserRepository) GetUser(userId typ.UserId) (ent.User, error) {
	u.lgr.LogFunctionInfo()
	ids := []typ.UserId{userId}
	users, err := u.GetUsers(ids)
	if err != nil {
		return ent.User{}, fmt.Errorf("failed to get user from database service: %w", err)
	}

	return users[0], nil
}

func (u *UserRepository) GetUsers(usrIds []typ.UserId) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()
	usrMs, err := u.dbS.GetUsers(usrIds)
	if err != nil {
		return nil, fmt.Errorf("failed to get user models from database service: %w", err)
	}

	return userEntitiesFromModels(usrMs), nil
}

func (u *UserRepository) NewUser(usrE ent.User) (ent.User, error) {
	u.lgr.LogFunctionInfo()

	usrM := userModelFromEntity(usrE)
	newUsrM, err := u.dbS.NewUser(usrM)
	if err != nil {
		return ent.User{}, fmt.Errorf("failed to create new user: %w", err)
	}

	return userEntityFromModel(newUsrM), nil
}

func userEntitiesFromModels(usrMs []model.User) []ent.User {
	var usrEs []ent.User
	for _, usrM := range usrMs {
		usrE := userEntityFromModel(usrM)
		usrEs = append(usrEs, usrE)
	}
	return usrEs
}

func userEntityFromModel(m model.User) ent.User {
	return ent.User{
		Id:      m.Id,
		Name:    m.Name,
		Email:   m.Email,
		PwdHash: m.PwdHash,
		Joined:  m.Joined,
	}
}

func userModelFromEntity(u ent.User) model.User {
	return model.User{
		Id:      u.Id,
		Name:    u.Name,
		Email:   u.Email,
		PwdHash: u.PwdHash,
		Joined:  u.Joined,
	}
}
