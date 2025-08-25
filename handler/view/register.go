package view

import (
	"net/http"
	mw "server/handler/middleware"
	"server/util"
	"text/template"
)

func Register() http.Handler {
	h := register{}
	return mw.AddMiddleware(h, mw.WithNoAuth(), mw.WithMethod(mw.GET))
}

type register struct{}

func (rh register) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	util.Log.FunctionInfo()

	tmpl, err := template.ParseFiles("./static/pages/register.html")
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

}
