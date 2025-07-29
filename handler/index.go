package handler

import (
	"errors"
	"log"
	"net/http"
	i "server/interfaces"
	"text/template"
)

func NewIndexHandler(l i.Logger, a i.AuthService) *IndexHandler {
	return &IndexHandler{
		lgr:   l,
		authS: a,
	}
}

type IndexHandler struct {
	authS i.AuthService
	lgr   i.Logger
}

func (i *IndexHandler) RenderIndexPage(w http.ResponseWriter, r *http.Request) {
	i.lgr.LogFunctionInfo()

	if r.Method != http.MethodGet {
		http.Error(w, "request method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var cookieFound bool
	cookie, err := r.Cookie("session_token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			cookieFound = true
		} else {
			cookieFound = false
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	if cookieFound {
		token := cookie.Value
		session, err := i.authS.ValidateAndRefreshSession(token)
		if err != nil {
			log.Println("error validating or refreshing session", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, session.Cookie())

		http.Redirect(w, r, "/chat", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("./static/pages/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
