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

func LeaveChat(a i.AuthService, c i.ChatService, m i.MessageService) http.Handler {
	h := leaveChat{
		chatS: c,
		msgS:  m,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.DELETE))
}

type leaveChat struct {
	chatS i.ChatService
	msgS  i.MessageService
}

func (h leaveChat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	query := r.URL.Query()
	req := lcrequest{
		ChatId: query.Get("ChatId"),
	}

	res, err := h.handleRequest(req, session.UserId())
	if err != nil {
		util.Log.Errorf("failed to handle leave chat request, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h leaveChat) handleRequest(req lcrequest, userId typ.UserId) (lcresponse, error) {
	util.Log.FunctionInfo()

	chatId, err := typ.ToChatId(req.ChatId)
	if err != nil {
		return lcresponse{}, err
	}

	chats, err := h.chatS.LeaveChat(chatId, userId)
	if err != nil {
		return lcresponse{}, err
	}

	newActiveChatId := chats[0].Id
	var messages []ent.Message
	if len(chats) != 0 {
		chatId := newActiveChatId
		messages, err = h.msgS.GetChatMessages(chatId, userId)
		if err != nil {
			return lcresponse{}, err
		}
	}

	newChatResponse := lcresponse{
		NewActiveChatId: newActiveChatId,
		Chats:           chats,
		Messages:        messages,
	}

	return newChatResponse, nil
}

type lcrequest struct {
	ChatId string `json:"ChatId"`
}

type lcresponse struct {
	Chats           []ent.Chat
	Messages        []ent.Message
	NewActiveChatId typ.ChatId
}
