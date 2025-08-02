package repository

import (
	"fmt"
	ent "server/data/entities"
	i "server/interfaces"
	cred "server/services/authService/credentials"
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

func (u *UserRepository) GetUsersForChat(chatId typ.ChatId) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()
	userMs, err := u.dbS.GetUsersForChat(chatId)
	if err != nil {
		return nil, err
	}

	return userEntitiesFromModels(userMs), nil
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

func (u *UserRepository) FindUserByEmail(email cred.Email) (ent.User, error) {
	u.lgr.LogFunctionInfo()

	userM, err := u.dbS.FindUser(email)
	if err != nil {
		return ent.User{}, err
	}

	if userM.Id == 0 {
		return ent.User{}, nil
	}

	userMs := []model.User{userM}
	userEs := userEntitiesFromModels(userMs)
	return userEs[0], nil
}

func (u *UserRepository) AddFriend(friend ent.Friend, userId typ.UserId) (ent.Friend, error) {
	u.lgr.LogFunctionInfo()

	friendM := model.Friend{
		UserAId: userId,
		UserBId: friend.Id,
	}

	friendM, err := u.dbS.InsertFriend(friendM)
	if err != nil {
		return ent.Friend{}, err
	}

	user, err := u.dbS.GetUser(userId)
	if err != nil {
		return ent.Friend{}, err
	}

	friendE := ent.Friend{
		Id:          user.Id,
		Name:        user.Name,
		Email:       user.Email,
		FriendSince: friendM.FriendSince,
	}

	return friendE, err
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
