package view

import (
	"net/http"
	mw "server/handler/middleware"
	"server/util"
	"text/template"
)

func Login() http.Handler {
	h := login{}
	return mw.AddMiddleware(h, mw.WithNoAuth(), mw.WithMethod(mw.GET))
}

type login struct{}

func (h login) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	tmpl, err := template.ParseFiles("./static/pages/login.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
