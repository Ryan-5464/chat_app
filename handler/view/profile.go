package view

import (
	"net/http"
	ent "server/data/entities"
	"server/handler/ctxutil"
	mw "server/handler/middleware"
	i "server/interfaces"
	ss "server/services/auth/session"
	"server/util"
	"text/template"
)

func Profile(a i.AuthService, u i.UserService) http.Handler {
	h := profile{
		userS: u,
	}
	return mw.AddMiddleware(h, mw.WithAuth(a), mw.WithMethod(mw.GET))
}

type profile struct {
	userS i.UserService
}

func (h profile) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	session := r.Context().Value(ctxutil.SessionKey).(ss.Session)

	user, err := h.userS.GetUser(session.UserId())
	if err != nil {
		util.Log.Errorf("failed to get user for user id %v", session.UserId())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := struct {
		User *ent.User
	}{
		User: user,
	}

	tmpl, err := template.ParseFiles("./static/pages/profile.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
