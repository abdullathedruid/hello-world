package handlers

import (
	"html/template"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles(
		"templates/layouts/base.html",
		"templates/pages/home.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
	}{
		Title: "HTMX + Go Demo",
	}

	tmpl.Execute(w, data)
}
