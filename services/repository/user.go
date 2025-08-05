package repository

import (
	"errors"
	"log"
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

func (u *UserRepository) NewUser(newUser ent.User) (*ent.User, error) {
	u.lgr.LogFunctionInfo()

	newUserModel := userModelFromEntity(newUser)
	userModel, err := u.dbS.NewUser(newUserModel)
	if err != nil {
		return nil, err
	}

	if userModel == nil {
		return nil, errors.New("new user missing!")
	}

	return userEntityFromModel(userModel), nil
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

	var contacts []ent.Contact

	contactRelations, err := u.dbS.GetContactRelations(userId)
	if err != nil {
		return contacts, err
	}

	if len(contactRelations) == 0 {
		return contacts, nil
	}

	contactIds := getContactIds(contactRelations)

	log.Println("CONTACTIDS:::::", contactIds)

	users, err := u.GetUsers(contactIds)
	if err != nil {
		return contacts, err
	}

	log.Println("USERS:::::", users)

	if len(users) == 0 {
		return contacts, errors.New("users missing!")
	}

	members, err := u.dbS.GetPrivateChatIdsForContacts(userId)
	if err != nil {
		return contacts, err
	}

	return createContacts(users, contactRelations, members), nil

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

func createContacts(users []ent.User, relations []model.ContactRelation, members []model.Member) []ent.Contact {
	var contacts []ent.Contact
	var userEmailMap = make(map[typ.UserId]cred.Email)
	var userNameMap = make(map[typ.UserId]string)
	var chatIdMap = make(map[typ.UserId]typ.ChatId)

	for _, user := range users {
		userEmailMap[user.Id] = user.Email
		userNameMap[user.Id] = user.Name
	}

	for _, member := range members {
		chatIdMap[member.UserId] = member.ChatId
	}

	for _, relation := range relations {
		contact := ent.Contact{
			Id:            relation.ContactId,
			KnownSince:    relation.Established,
			Name:          userNameMap[relation.ContactId],
			Email:         userEmailMap[relation.ContactId],
			PrivateChatId: chatIdMap[relation.ContactId],
		}
		contacts = append(contacts, contact)
	}
	return contacts
}

func getContactIds(c []model.ContactRelation) []typ.UserId {
	var contactIds []typ.UserId
	for _, relation := range c {
		contactIds = append(contactIds, relation.ContactId)
	}
	return contactIds
}
