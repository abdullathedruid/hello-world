package handlers

import (
	"html/template"
	"net/http"

	"github.com/sirupsen/logrus"
)

func DebugHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Debug page accessed")
	tmpl, err := template.ParseFiles(
		"templates/layouts/base.html",
		"templates/pages/debug.html",
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	data := struct {
		Title string
	}{
		Title: "Farcaster MiniApp Debug",
	}
	
	tmpl.Execute(w, data)
}