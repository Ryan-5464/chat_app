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

func SwitchChat(a i.AuthService, m i.MessageService) http.Handler {
	h := switchChat{
		msgS: m,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.GET))
}

type switchChat struct {
	msgS i.MessageService
}

func (h switchChat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	query := r.URL.Query()
	req := screquest{
		ChatId: query.Get("ChatId"),
	}

	res, err := h.handleRequest(req, session.UserId())
	if err != nil {
		util.Log.Errorf("failed to handle switch chat screquest, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h switchChat) handleRequest(req screquest, userId typ.UserId) (scresponse, error) {
	util.Log.FunctionInfo()

	chatId, err := typ.ToChatId(req.ChatId)
	if err != nil {
		return scresponse{}, err
	}

	messages, err := h.msgS.GetChatMessages(chatId, userId)
	if err != nil {
		return scresponse{}, err
	}

	latestMsgId := findLastestMessageId(messages)

	if err := h.msgS.UpdateLastReadMsgId(latestMsgId, chatId, userId); err != nil {
		return scresponse{}, err
	}

	res := scresponse{
		ActiveChatId: chatId,
		Messages:     messages,
	}

	return res, nil
}

type screquest struct {
	ChatId string `json:"ChatId"`
}

type scresponse struct {
	ActiveChatId typ.ChatId    `json:"ActiveChatId"`
	Messages     []ent.Message `json:"Messages"`
}
