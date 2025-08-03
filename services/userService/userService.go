package userservice

import (
	"fmt"
	ent "server/data/entities"
	i "server/interfaces"
	cred "server/services/authService/credentials"
	model "server/services/dbService/SQL/models"
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

func (u *UserService) AddContact(contactEmail cred.Email, userId typ.UserId) ([]ent.Contact, error) {
	u.lgr.LogFunctionInfo()

	var contacts []ent.Contact
	var contactEmails []cred.Email
	contactEmails = append(contactEmails, contactEmail)
	users, err := u.usrR.FindUsers(contactEmails)
	if err != nil {
		return contacts, err
	}

	if len(users) == 0 {
		return contacts, err
	}

	contactRelation := model.ContactRelation{
		UserId:    userId,
		ContactId: users[0].Id,
	}

	contactRelations, err := u.usrR.AddContact(contactRelation)
	if err != nil {
		return contacts, err
	}

	contact := ent.Contact{
		Id:         users[0].Id,
		Name:       users[0].Name,
		Email:      contactEmail,
		KnownSince: contactRelations[0].Established,
	}

	contacts = append(contacts, contact)
	return contacts, nil
}

func (u *UserService) GetContacts(userId typ.UserId) ([]ent.Contact, error) {
	u.lgr.LogFunctionInfo()

	var contacts []ent.Contact

	contactRelations, err := u.usrR.GetContactRelations(userId)
	if err != nil {
		return contacts, err
	}

	if len(contactRelations) == 0 {
		return contacts, nil
	}

	contactIds := getContactIds(contactRelations)

	users, err := u.usrR.GetUsers(contactIds)
	if err != nil {
		return contacts, err
	}

	return createContacts(users, contactRelations), nil
}

func createContacts(users []model.User, crs []model.ContactRelation) []ent.Contact {
	var contacts []ent.Contact
	for i := 0; i < len(users); i++ {
		contact := ent.Contact{
			Id:         crs[i].Contact2,
			KnownSince: crs[i].Established,
		}
		// Can probably optimize this with a better data structure
		for _, user := range users {
			if user.Id == crs[i].Contact1 {
				contact.Name = user.Name
				contact.Email = user.Email
				break
			}
		}
		contacts = append(contacts, contact)
	}
	return contacts
}

func getContactIds(c []model.ContactRelation) []typ.UserId {
	var contactIds []typ.UserId
	for _, relation := range c {
		contactIds = append(contactIds, relation.Contact2)
	}
	return contactIds
}
