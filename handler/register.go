package handler

import (
	"net/http"
	i "server/interfaces"
	"text/template"
)

func NewRegisterHandler(l i.Logger) *RegisterHandler {
	return &RegisterHandler{
		lgr: l,
	}
}

type RegisterHandler struct {
	lgr i.Logger
}

func (rh *RegisterHandler) RenderRegisterPage(w http.ResponseWriter, r *http.Request) {
	rh.lgr.LogFunctionInfo()

	tmpl, err := template.ParseFiles("./static/pages/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}

func (rh *RegisterHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	rh.lgr.LogFunctionInfo()

}
