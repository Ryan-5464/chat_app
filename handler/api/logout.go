package api

import (
	"net/http"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	"server/util"
)

func Logout(a i.AuthService, cn i.ConnectionService) http.Handler {
	h := logout{
		connS: cn,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.GET))
}

type logout struct {
	connS i.ConnectionService
}

func (h logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	if err := h.connS.ChangeOnlineStatus("offline", session.UserId()); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	session.InvalidateSession()

	http.SetCookie(w, session.Cookie())

	util.Log.Dbug("token invalidated <<<>>>>")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
