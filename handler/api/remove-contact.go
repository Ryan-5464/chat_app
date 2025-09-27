package api

import (
	"net/http"
	ent "server/data/entities"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	typ "server/types"
	"server/util"
)

func RemoveContact(a i.AuthService, c i.ChatService, m i.MessageService, u i.UserService, cn i.ConnectionService) http.Handler {
	h := removeContact{
		chatS: c,
		msgS:  m,
		userS: u,
		connS: cn,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.DELETE))
}

type removeContact struct {
	chatS i.ChatService
	msgS  i.MessageService
	userS i.UserService
	connS i.ConnectionService
}

func (h removeContact) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	query := r.URL.Query()
	req := rcrequest{
		ContactId: query.Get("ContactId"),
	}

	res, err := h.handleRequest(req, session.UserId())
	if err != nil {
		util.Log.Errorf("failed to handle remove contact request, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h removeContact) handleRequest(req rcrequest, userId typ.UserId) (rcresponse, error) {
	util.Log.FunctionInfo()

	contactId, err := typ.ToContactId(req.ContactId)
	if err != nil {
		return rcresponse{}, err
	}

	if err := h.userS.RemoveContact(contactId, userId); err != nil {
		return rcresponse{}, err
	}

	chats, err := h.chatS.GetChats(userId)
	if err != nil {
		return rcresponse{}, err
	}

	var newActiveChatId typ.ChatId
	var messages []ent.Message
	if len(chats) != 0 {
		newActiveChatId = chats[0].Id
		messages, err = h.msgS.GetChatMessages(newActiveChatId, userId)
		if err != nil {
			return rcresponse{}, err
		}
	} else {
		newActiveChatId = typ.ChatId(0)
		messages = []ent.Message{}
	}

	contacts, err := h.userS.GetContacts(typ.UserId(contactId))
	if err != nil {
		return rcresponse{}, err
	}

	conn := h.connS.GetConnection(typ.UserId(contactId))

	payload := struct {
		Type     string
		Contacts []ent.Contact
	}{
		Type:     "RemoveContact",
		Contacts: contacts,
	}

	if err := conn.WriteJSON(payload); err != nil {
		util.Log.Errorf("failed to write to websocket connection: %v", err)
		return rcresponse{}, err
	}

	newChatResponse := rcresponse{
		NewActiveChatId: newActiveChatId,
		Contacts:        contacts,
		Messages:        messages,
	}

	return newChatResponse, nil
}

func (h removeContact) GetContactsForRemovedContact(contactId typ.UserId) ([]ent.Contact, error) {
	util.Log.FunctionInfo()

	contacts, err := h.userS.GetContacts(contactId)
	if err != nil {
		return nil, err
	}

	return contacts, nil
}

type rcrequest struct {
	ContactId string `json:"ContactId"`
}

type rcresponse struct {
	Contacts        []ent.Contact
	Messages        []ent.Message
	NewActiveChatId typ.ChatId
}
