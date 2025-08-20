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

func (u *UserRepository) GetUser(userId typ.UserId) (*ent.User, error) {
	u.lgr.LogFunctionInfo()

	userModel, err := u.dbS.GetUser(userId)
	if err != nil {
		return nil, err
	}

	return userEntityFromModel(userModel), nil

}

func (u *UserRepository) GetUsers(userIds []typ.UserId) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()

	userModels, err := u.dbS.GetUsers(userIds)
	if err != nil {
		return []ent.User{}, err
	}

	return userEntitiesFromModels(userModels), nil

}

func (u *UserRepository) GetChatUsers(chatId typ.ChatId) ([]ent.User, error) {
	u.lgr.LogFunctionInfo()

	members, err := u.dbS.GetMembers(chatId)
	if err != nil {
		return []ent.User{}, err
	}

	memberIds := getMemberIds(members)

	users, err := u.dbS.GetUsers(memberIds)
	if err != nil {
		return []ent.User{}, err
	}

	return userEntitiesFromModels(users), nil
}

func (u *UserRepository) NewUser(userName string, userEmail cred.Email, pwdHash cred.PwdHash) (*ent.User, error) {
	u.lgr.LogFunctionInfo()

	lastInsertId, err := u.dbS.CreateUser(userName, userEmail, pwdHash)
	if err != nil {
		return nil, err
	}

	user, err := u.dbS.GetUser(typ.UserId(lastInsertId))
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

	userModels, err := u.dbS.FindUsers(emails)
	if err != nil {
		return []ent.User{}, err
	}

	return userEntitiesFromModels(userModels), nil
}

func (u *UserRepository) AddContact(contactId typ.ContactId, contactName string, contactEmail cred.Email, userId typ.UserId) (*ent.Contact, error) {
	u.lgr.LogFunctionInfo()

	lastInsertId, err := u.dbS.CreateContact(userId, typ.ContactId(contactId))
	if err != nil {
		return nil, err
	}

	contactModel, err := u.dbS.GetContact(typ.ChatId(lastInsertId))
	if err != nil {
		return nil, err
	}

	if contactModel == nil {
		return nil, errors.New("added contact relation missing!")
	}

	contactEnt := &ent.Contact{
		Id:            decipherContactId(contactModel.Id1, contactModel.Id2, userId),
		Name:          contactName,
		Email:         contactEmail,
		KnownSince:    contactModel.CreatedAt,
		ContactChatId: contactModel.ChatId,
	}

	return contactEnt, nil

}

func (u *UserRepository) GetContacts(userId typ.UserId) ([]ent.Contact, error) {
	u.lgr.LogFunctionInfo()

	contactModels, err := u.dbS.GetContacts(userId)
	if err != nil {
		return []ent.Contact{}, err
	}

	contacts := createContacts(contactModels, userId)

	var contactIds []typ.UserId
	for _, contact := range contacts {
		contactIds = append(contactIds, typ.UserId(contact.Id))
	}

	cs, err := u.dbS.GetUsers(contactIds)
	if err != nil {
		return []ent.Contact{}, err
	}

	return mapUserInfoToContacts(contacts, cs), nil
}

func (u *UserRepository) GetContact(chatId typ.ChatId, userId typ.UserId) (*ent.Contact, error) {
	u.lgr.LogFunctionInfo()
	contactModel, err := u.dbS.GetContact(chatId)
	if err != nil {
		return nil, err
	}

	contact := createContact(contactModel, userId)

	user, err := u.dbS.GetUser(typ.UserId(contact.Id))
	if err != nil {
		return nil, err
	}

	contact.Name = user.Name
	contact.Email = user.Email

	return contact, nil
}

func (u *UserRepository) RemoveContact(contactId typ.ContactId, userId typ.UserId) error {
	u.lgr.LogFunctionInfo()
	return u.dbS.DeleteContact(contactId, userId)
}

func (u *UserRepository) EditUserName(name string, userId typ.UserId) error {
	u.lgr.LogFunctionInfo()
	return u.dbS.UpdateUserName(name, userId)
}

func mapUserInfoToContacts(contacts []ent.Contact, users []model.User) []ent.Contact {
	contactIdMap := make(map[typ.ContactId]ent.Contact)
	for _, contact := range contacts {
		contactIdMap[contact.Id] = contact
	}

	updatedContacts := []ent.Contact{}
	for _, user := range users {
		contact := contactIdMap[typ.ContactId(user.Id)]
		contact.Name = user.Name
		contact.Email = user.Email
		updatedContacts = append(updatedContacts, contact)
	}

	return updatedContacts
}

func createContact(c *model.Contact, userId typ.UserId) *ent.Contact {
	return &ent.Contact{
		Id:            decipherContactId(c.Id1, c.Id2, userId),
		ContactChatId: c.ChatId,
		KnownSince:    c.CreatedAt,
	}
}

func createContacts(contactModels []model.Contact, userId typ.UserId) []ent.Contact {
	var contacts []ent.Contact
	for _, model := range contactModels {

		contact := ent.Contact{
			Id:            decipherContactId(model.Id1, model.Id2, userId),
			ContactChatId: model.ChatId,
			KnownSince:    model.CreatedAt,
		}

		contacts = append(contacts, contact)
	}

	return contacts
}

func decipherContactId(id1 typ.UserId, id2 typ.UserId, userId typ.UserId) typ.ContactId {
	if id1 == userId {
		return typ.ContactId(id2)
	} else {
		return typ.ContactId(id1)
	}
}

func userEntitiesFromModels(userModels []model.User) []ent.User {
	if len(userModels) == 0 {
		return []ent.User{}
	}

	var usrEs []ent.User
	for _, m := range userModels {
		usrE := ent.User{
			Id:      m.Id,
			Name:    m.Name,
			Email:   m.Email,
			PwdHash: m.PwdHash,
			Joined:  m.Joined,
		}
		usrEs = append(usrEs, usrE)
	}
	return usrEs
}

func userEntityFromModel(m *model.User) *ent.User {
	if m == nil {
		return nil
	}

	return &ent.User{
		Id:      m.Id,
		Name:    m.Name,
		Email:   m.Email,
		PwdHash: m.PwdHash,
		Joined:  m.Joined,
	}
}

func getMemberIds(members []model.Member) []typ.UserId {
	var userIds []typ.UserId
	for _, member := range members {
		userIds = append(userIds, member.UserId)
	}
	return userIds
}
