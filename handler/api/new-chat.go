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

func NewChat(a i.AuthService, c i.ChatService, m i.MessageService) http.Handler {
	h := newChat{
		chatS: c,
		msgS:  m,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.POST))
}

type newChat struct {
	chatS i.ChatService
	msgS  i.MessageService
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

	var replyId typ.MessageId
	if err := h.msgS.HandleNewMessage(userId, chat.Id, &replyId, "Add someone to chat with!"); err != nil {
		return ncresponse{}, err
	}

	messages, err := h.msgS.GetChatMessages(chat.Id, userId)
	if err != nil {
		return ncresponse{}, err
	}

	res := ncresponse{
		Chats:        []ent.Chat{*chat},
		Messages:     messages,
		ActiveChatId: chat.Id,
	}

	return res, nil
}

type ncrequest struct {
	Name string `json:"Name"`
}

type ncresponse struct {
	Chats        []ent.Chat
	ActiveChatId typ.ChatId    `json:"ActiveChatId"`
	Messages     []ent.Message `json:"Messages"`
}
