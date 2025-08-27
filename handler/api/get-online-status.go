package api

import (
	"net/http"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	"server/util"
)

func GetOnlineStatus(a i.AuthService, cn i.ConnectionService) http.Handler {
	h := getOnlineStatus{
		connS: cn,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.GET))
}

type getOnlineStatus struct {
	connS i.ConnectionService
}

func (h getOnlineStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	status := h.connS.GetOnlineStatus(session.UserId())

	res := gosresponse{
		OnlineStatus: status,
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("response sent")
}

type gosresponse struct {
	OnlineStatus string
}
