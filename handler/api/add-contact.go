package api

import (
	"encoding/json"
	"errors"
	"net/http"
	ent "server/data/entities"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	cred "server/services/auth/credentials"
	ss "server/services/auth/session"
	typ "server/types"
	"server/util"
)

func AddContact(a i.AuthService, cn i.ConnectionService, u i.UserService) http.Handler {
	h := addContact{
		authS: a,
		connS: cn,
		userS: u,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.POST))
}

type addContact struct {
	authS i.AuthService
	connS i.ConnectionService
	userS i.UserService
}

func (h addContact) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	var req acrequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Log.Errorf("failed to decode JSON request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	util.Log.Dbugf("Add contact request name: %v", req.Email)

	res, err := h.handleRequest(req, session.UserId())
	if err != nil {
		util.Log.Errorf("failed to handle add contact request, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h addContact) handleRequest(req acrequest, userId typ.UserId) (acresponse, error) {
	util.Log.FunctionInfo()

	var res acresponse

	contact, err := h.userS.AddContact(cred.Email(req.Email), userId)
	if err != nil {
		return res, err
	}

	if contact == nil {
		return res, errors.New("failed to add contact")
	}

	contacts, err := h.GetContactsForNewContact(typ.UserId(contact.Id))
	if err != nil {
		return res, err
	}

	conn := h.connS.GetConnection(typ.UserId(contact.Id))

	payload := struct {
		Type     string
		Contacts []ent.Contact
	}{
		Type:     "AddContact",
		Contacts: contacts,
	}

	if err := conn.WriteJSON(payload); err != nil {
		util.Log.Errorf("failed to write to websocket connection: %v", err)
		return acresponse{}, err
	}

	res = acresponse{
		Contacts: []ent.Contact{*contact},
	}

	return res, nil
}

func (h addContact) GetContactsForNewContact(contactId typ.UserId) ([]ent.Contact, error) {
	util.Log.FunctionInfo()

	contacts, err := h.userS.GetContacts(contactId)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

type acrequest struct {
	Email string `json:"Email"`
}

type acresponse struct {
	Contacts []ent.Contact
}
