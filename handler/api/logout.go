package api

import (
	"net/http"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	"server/util"
)

func Logout(a i.AuthService) http.Handler {
	h := logout{}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.GET))
}

type logout struct {
}

func (h logout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	session.InvalidateSession()

	http.SetCookie(w, session.Cookie())

	util.Log.Dbug("token invalidated <<<>>>>")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
