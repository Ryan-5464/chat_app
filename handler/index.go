package handler

import (
	"net/http"
	i "server/interfaces"
	ss "server/services/authService/session"
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

	session := r.Context().Value("session").(ss.Session)
	emptySession := ss.Session{}
	if session != emptySession {
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
