package api

import (
	"net/http"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	typ "server/types"
	"server/util"
)

func DeleteMessage(a i.AuthService, s i.MessageService) http.Handler {
	h := deleteMessage{
		msgS: s,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.DELETE))
}

type deleteMessage struct {
	msgS i.MessageService
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
	}

	res, err := h.handleRequest(req, session.UserId())
	if err != nil {
		util.Log.Errorf("failed to delete message: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h deleteMessage) handleRequest(req dmrequest, userId typ.UserId) (dmresponse, error) {
	util.Log.FunctionInfo()

	messageId, err := typ.ToMessageId(req.MessageId)
	if err != nil {
		return dmresponse{}, err
	}

	if err := h.msgS.DeleteMessage(messageId); err != nil {
		return dmresponse{}, err
	}

	return dmresponse{
		MessageId: messageId,
	}, nil
}

type dmrequest struct {
	MessageId string `json:"MessageId"`
	UserId    string `json:"UserId"`
}

type dmresponse struct {
	MessageId typ.MessageId
}
