package handler

import (
	"net/http"
	"text/template"
)

func NewIndexHandler() *IndexHandler {
	return &IndexHandler{}
}

type IndexHandler struct {
}

func (i *IndexHandler) RenderIndexPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/pages/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
