package api

import (
	"encoding/json"
	"net/http"
	ent "server/data/entities"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	typ "server/types"
	"server/util"
)

func NewChat(a i.AuthService, c i.ChatService) http.Handler {
	h := newChat{
		chatS: c,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.POST))
}

type newChat struct {
	chatS i.ChatService
}

func (h newChat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	var req ncrequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Log.Errorf("failed to decode JSON request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	util.Log.Dbugf("New chat request name: %v", req.Name)

	res, err := h.handleNewChatRequest(req, session.UserId())
	if err != nil {
		util.Log.Errorf("failed to handle new chat request, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h newChat) handleNewChatRequest(cr ncrequest, userId typ.UserId) (ncresponse, error) {
	util.Log.FunctionInfo()

	chat, err := h.chatS.NewChat(cr.Name, userId)
	if err != nil {
		return ncresponse{}, err
	}

	res := ncresponse{
		Chats: []ent.Chat{*chat},
	}

	return res, nil
}

type ncrequest struct {
	Name string `json:"Name"`
}

type ncresponse struct {
	Chats []ent.Chat
}
