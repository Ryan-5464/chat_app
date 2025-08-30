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

func RemoveMember(a i.AuthService, c i.ChatService, cn i.ConnectionService) http.Handler {
	h := removeMember{
		chatS: c,
		connS: cn,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.DELETE))

}

type removeMember struct {
	chatS i.ChatService
	connS i.ConnectionService
}

func (h removeMember) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	query := r.URL.Query()
	req := rmrequest{
		ChatId: query.Get("ChatId"),
		UserId: query.Get("UserId"),
	}

	res, err := h.handleRequest(req, session.UserId())
	if err != nil {
		util.Log.Errorf("failed to handle remove member request, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

func (h removeMember) handleRequest(req rmrequest, adminId typ.UserId) (rmresponse, error) {
	util.Log.FunctionInfo()

	chatId, err := typ.ToChatId(req.ChatId)
	if err != nil {
		return rmresponse{}, err
	}

	memberId, err := typ.ToUserId(req.UserId)
	if err != nil {
		return rmresponse{}, err
	}

	if err := h.chatS.RemoveMember(chatId, memberId, adminId); err != nil {
		return rmresponse{}, err
	}

	conn := h.connS.GetConnection(memberId)

	payload := struct {
		Type   string
		ChatId typ.ChatId
	}{
		Type:   "RemoveMember",
		ChatId: chatId,
	}

	if err := conn.WriteJSON(payload); err != nil {
		util.Log.Errorf("failed to write to websocket connection: %v", err)
		return rmresponse{}, err
	}

	return rmresponse{Success: true}, nil
}

type rmrequest struct {
	ChatId string `json:"ChatId"`
	UserId string `json:"UserId"`
}

type rmresponse struct {
	Success bool
}
