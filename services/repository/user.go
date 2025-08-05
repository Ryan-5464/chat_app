package repository

import (
	"errors"
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

func (u *UserRepository) GetUsers(userIds []typ.UserId) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()

	var users []ent.User

	userModels, err := u.dbS.GetUsers(userIds)
	if err != nil {
		return users, err
	}

	if len(userModels) == 0 {
		return users, nil
	}

	return userEntitiesFromModels(userModels), nil

}

func (u *UserRepository) GetChatUsers(chatId typ.ChatId) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()

	var users []ent.User

	userModels, err := u.dbS.GetChatUsers(chatId)
	if err != nil {
		return users, err
	}

	if len(userModels) == 0 {
		return users, nil
	}

	return userEntitiesFromModels(userModels), nil
}

func (u *UserRepository) NewUser(userName string, userEmail cred.Email, pwdHash cred.PwdHash) (*ent.User, error) {
	u.lgr.LogFunctionInfo()

	lastInsertId, err := u.dbS.CreateUser(userName, userEmail, pwdHash)
	if err != nil {
		return nil, err
	}

	user, err := u.dbS.GetNewUser(lastInsertId)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("new user missing!")
	}

	return userEntityFromModel(user), nil
}

func (u *UserRepository) FindUser(email cred.Email) (*model.User, error) {
	u.lgr.LogFunctionInfo()
	return u.dbS.FindUser(email)
}

func (u *UserRepository) FindUsers(emails []cred.Email) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()

	var users []ent.User

	userModels, err := u.dbS.FindUsers(emails)
	if err != nil {
		return users, err
	}

	if len(userModels) == 0 {
		return users, nil
	}

	return userEntitiesFromModels(userModels), nil
}

func (u *UserRepository) AddContact(contact ent.Contact, userId typ.UserId) (*ent.Contact, error) {
	u.lgr.LogFunctionInfo()

	contactRelation, err := u.dbS.AddContactRelation(userId, contact.Id)
	if err != nil {
		return nil, err
	}

	if contactRelation == nil {
		return nil, errors.New("added contact relation missing!")
	}

	c := &ent.Contact{
		Id:         contact.Id,
		Name:       contact.Name,
		Email:      contact.Email,
		KnownSince: contactRelation.Established,
	}

	return c, nil

}

func (u *UserRepository) GetContacts(userId typ.UserId) ([]ent.Contact, error) {
	u.lgr.LogFunctionInfo()

	contactChats, err := u.dbS.GetContactChats(userId)
	if err != nil {
		return []ent.Contact{}, err
	}

	contacts, err := createContacts(contactChats, userId)
	if err != nil {
		return []ent.Contact{}, err
	}

	var contactIds []typ.UserId
	for _, contact := range contacts {
		contactIds = append(contactIds, contact.Id)
	}

	users, err := u.dbS.GetUsers(contactIds)
	if err != nil {
		return []ent.Contact{}, err
	}

	return mapUserInfoToContacts(contacts, users), nil
}

func mapUserInfoToContacts(contacts []ent.Contact, users []model.User) []ent.Contact {
	contactIdMap := make(map[typ.UserId]ent.Contact)
	for _, contact := range contacts {
		contactIdMap[contact.Id] = contact
	}

	updatedContacts := []ent.Contact{}
	for _, user := range users {
		contact := contactIdMap[user.Id]
		contact.Name = user.Name
		contact.Email = user.Email
		updatedContacts = append(updatedContacts, contact)
	}

	return updatedContacts
}

func createContacts(chats []model.ContactChat, userId typ.UserId) ([]ent.Contact, error) {
	var contacts []ent.Contact
	for _, chat := range chats {

		contactId, err := decipherContactId(chat, userId)
		if err != nil {
			return []ent.Contact{}, err
		}

		contact := ent.Contact{
			Id:            contactId,
			ContactChatId: chat.Id,
		}

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func decipherContactId(chat model.ContactChat, userId typ.UserId) (typ.UserId, error) {
	var contactId typ.UserId
	if chat.Member1Id == userId {
		contactId = chat.Member2Id
	}
	if chat.Member2Id == userId {
		contactId = chat.Member1Id
	} else {
		return 0, errors.New("no matching id for contact")
	}
	return contactId, nil
}

func userEntitiesFromModels(usrMs []model.User) []ent.User {
	var usrEs []ent.User
	for _, usrM := range usrMs {
		usrE := userEntityFromModel(&usrM)
		usrEs = append(usrEs, *usrE)
	}
	return usrEs
}

func userEntityFromModel(m *model.User) *ent.User {
	return &ent.User{
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
