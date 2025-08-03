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

func (u *UserRepository) GetUsers(usrIds []typ.UserId) ([]model.User, error) {
	u.lgr.LogFunctionInfo()
	return u.dbS.GetUsers(usrIds)
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

func (u *UserRepository) FindUsers(emails []cred.Email) ([]model.User, error) {
	u.lgr.LogFunctionInfo()
	return u.dbS.FindUsers(emails)
}

func (u *UserRepository) AddContact(c model.ContactRelation) ([]model.ContactRelation, error) {
	u.lgr.LogFunctionInfo()

	var contactRelations []model.ContactRelation

	contactRelations, err := u.dbS.AddContactRelation(c)

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

func (u *UserRepository) GetContactRelations(userId typ.UserId) ([]model.ContactRelation, error) {
	u.lgr.LogFunctionInfo()

	return u.dbS.GetContactRelations(userId)
}

func (u *UserRepository) GetFriends(userId typ.UserId) ([]ent.Friend, error) {
	u.lgr.LogFunctionInfo()

	friendMs, err := u.dbS.GetFriends(userId)
	if err != nil {
		return []ent.Friend{}, err
	}

	if len(friendMs) == 0 {
		return []ent.Friend{}, nil
	}

	return friendEntitiesFromModels(friendMs), nil
}

func friendEntitiesFromModels(friendMs []model.Friend) []ent.Friend {
	var friendEs []ent.Friend
	for _, friendM := range friendMs {
		friendE := ent.Friend{
			Id:          friendM.UserAId,
			FriendSince: friendM.FriendSince,
		}
		friendEs = append(friendEs, friendE)
	}
	return friendEs
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
