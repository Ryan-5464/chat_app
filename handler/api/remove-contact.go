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

func RemoveContact(a i.AuthService, c i.ChatService, m i.MessageService, u i.UserService) http.Handler {
	h := removeContact{
		chatS: c,
		msgS:  m,
		userS: u,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.DELETE))
}

type removeContact struct {
	chatS i.ChatService
	msgS  i.MessageService
	userS i.UserService
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

	newActiveChatId := chats[0].Id
	var messages []ent.Message
	if len(chats) != 0 {
		chatId := newActiveChatId
		messages, err = h.msgS.GetChatMessages(chatId, userId)
		if err != nil {
			return rcresponse{}, err
		}
	}

	contacts, err := h.userS.GetContacts(userId)
	if err != nil {
		return rcresponse{}, err
	}

	newChatResponse := rcresponse{
		NewActiveChatId: newActiveChatId,
		Contacts:        contacts,
		Messages:        messages,
	}

	return newChatResponse, nil
}

type rcrequest struct {
	ContactId string `json:"ContactId"`
}

type rcresponse struct {
	Contacts        []ent.Contact
	Messages        []ent.Message
	NewActiveChatId typ.ChatId
}
