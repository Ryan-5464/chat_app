package handler

import (
	"net/http"
	i "server/interfaces"
	"text/template"
)

func NewLoginHandler(l i.Logger) *LoginHandler {
	return &LoginHandler{
		lgr: l,
	}
}

type LoginHandler struct {
	lgr i.Logger
}

func (l *LoginHandler) RenderLoginPage(w http.ResponseWriter, r *http.Request) {
	l.lgr.LogFunctionInfo()

	tmpl, err := template.ParseFiles("./static/pages/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
