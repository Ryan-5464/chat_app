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

func SwitchContactChat(a i.AuthService, m i.MessageService) http.Handler {
	h := switchContactChat{
		msgS: m,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.GET))
}

type switchContactChat struct {
	msgS i.MessageService
}

func (h switchContactChat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	query := r.URL.Query()
	req := sccrequest{
		ContactChatId: query.Get("ContactChatId"),
	}

	res, err := h.handleRequest(req, session.UserId())
	if err != nil {
		util.Log.Errorf("failed to handle contact chat switch request, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h switchContactChat) handleRequest(s sccrequest, userId typ.UserId) (sccresponse, error) {
	util.Log.FunctionInfo()

	contactChatId, err := typ.ToChatId(s.ContactChatId)
	if err != nil {
		return sccresponse{}, err
	}

	messages, err := h.msgS.GetContactMessages(contactChatId, userId)
	if err != nil {
		return sccresponse{}, err
	}

	res := sccresponse{
		ActiveContactChatId: contactChatId,
		Messages:            messages,
	}

	return res, nil
}

type sccrequest struct {
	ContactChatId string `json:"ContactChatId"`
}

type sccresponse struct {
	ActiveContactChatId typ.ChatId
	Messages            []ent.Message
}
