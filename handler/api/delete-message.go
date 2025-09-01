package api

import (
	"log"
	"net/http"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	typ "server/types"
	"server/util"
)

func DeleteMessage(a i.AuthService, m i.MessageService, c i.ChatService, cn i.ConnectionService) http.Handler {
	h := deleteMessage{
		msgS:  m,
		chatS: c,
		connS: cn,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.DELETE))
}

type deleteMessage struct {
	msgS  i.MessageService
	chatS i.ChatService
	connS i.ConnectionService
}

func (h deleteMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	query := r.URL.Query()
	userId, err := typ.ToUserId(query.Get("UserId"))
	if err != nil {
		util.Log.Errorf("failed to parse userId: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if session.UserId() != typ.UserId(userId) {
		util.Log.Info("error: user unauthorized to delete message")
		return
	}

	req := dmrequest{
		MessageId: query.Get("MessageId"),
		UserId:    query.Get("UserId"),
		ChatId:    query.Get("ChatId"),
	}

	res, err := h.handleRequest(req)
	if err != nil {
		util.Log.Errorf("failed to delete message: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h deleteMessage) handleRequest(req dmrequest) (dmresponse, error) {
	util.Log.FunctionInfo()

	messageId, err := typ.ToMessageId(req.MessageId)
	if err != nil {
		return dmresponse{}, err
	}

	if err := h.msgS.DeleteMessage(messageId); err != nil {
		return dmresponse{}, err
	}

	chatId, err := typ.ToChatId(req.ChatId)
	if err != nil {
		return dmresponse{}, err
	}

	members, err := h.chatS.GetChatMembers(chatId)
	if err != nil {
		return dmresponse{}, err
	}

	for _, member := range members {
		conn := h.connS.GetConnection(member.UserId)
		if conn == nil {
			continue
		}

		res := dmresponse{
			Type:      "DeleteMessage",
			MessageId: messageId,
		}

		log.Println("DELETED MESSAGE", res)

		if err := conn.WriteJSON(res); err != nil {
			util.Log.Errorf("failed to write to websocket connection: %v", err)
			return dmresponse{}, err
		}
	}

	return dmresponse{
		Type:      "DeleteMessage",
		MessageId: messageId,
	}, nil
}

type dmrequest struct {
	MessageId string `json:"MessageId"`
	UserId    string `json:"UserId"`
	ChatId    string `json:"ChatId"`
}

type dmresponse struct {
	Type      string
	MessageId typ.MessageId
}
