package api

import (
	"net/http"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	"server/util"
)

func OnlineStatus(a i.AuthService, cn i.ConnectionService) http.Handler {
	h := onlineStatus{
		connS: cn,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.GET))
}

type onlineStatus struct {
	connS i.ConnectionService
}

func (o onlineStatus) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	query := r.URL.Query()

	if err := o.connS.ChangeOnlineStatus(query.Get("Status"), session.UserId()); err != nil {
		util.Log.Errorf("failed to get contacts to broadcast online status, Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	res := osresponse{
		Status: query.Get("Status"),
	}

	SendJSONResponse(w, res)

	util.Log.Dbug("->>>> RESPONSE SENT")
}

type osresponse struct {
	Status string
}
