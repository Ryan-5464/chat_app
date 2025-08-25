package api

import (
	"encoding/json"
	"net/http"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	typ "server/types"
	"server/util"
)

func EditMessage(a i.AuthService, m i.MessageService) http.Handler {
	h := editMessage{
		msgS: m,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.POST))
}

type editMessage struct {
	msgS i.MessageService
}

func (h editMessage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	var req emrequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.Log.Errorf("failed to decode JSON request body: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userId, err := typ.ToUserId(req.UserId)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if session.UserId() != typ.UserId(userId) {
		util.Log.Errorf("user does not own message: %v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	res, err := h.handleRequest(req)
	if err != nil {
		util.Log.Errorf("failed to handle contact edit chat name request, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbugf("->>>> RESPONSE SENT:: %v", res)

}

func (h editMessage) handleRequest(req emrequest) (emresponse, error) {
	util.Log.FunctionInfo()

	messageId, err := typ.ToMessageId(req.MessageId)
	if err != nil {
		return emresponse{}, err
	}

	msg, err := h.msgS.EditMessage(req.MsgText, messageId)
	if err != nil {
		return emresponse{}, err
	}

	return emresponse{
		MsgText: msg.Text,
	}, nil

}

type emrequest struct {
	MsgText   string `json:"MsgText"`
	MessageId string `json:"MessageId"`
	UserId    string `json:"UserId"`
}

type emresponse struct {
	MsgText string
}
