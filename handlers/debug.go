package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
)

func DebugHandler(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "Debug page accessed")
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
