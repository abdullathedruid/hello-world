package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"

	"hello-world/services"
)

var clickService = services.NewClickService()

// ResetClickService resets the click counter for testing
func ResetClickService() {
	clickService.Reset()
}

func TimeFragmentHandler(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "Time endpoint accessed")
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	tmpl := `<div class="bg-blue-50 border border-blue-200 rounded-md p-4">
		<p class="text-gray-700">Current time: <strong class="text-blue-600">{{.Time}}</strong></p>
		<button class="mt-3 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
		        hx-get="/api/time" hx-target="#time-display">Refresh Time</button>
	</div>`

	t, _ := template.New("time").Parse(tmpl)
	data := struct{ Time string }{Time: currentTime}
	t.Execute(w, data)
}

func ClickFragmentHandler(w http.ResponseWriter, r *http.Request) {
	count := clickService.IncrementClick()
	slog.InfoContext(r.Context(), "Button clicked", "count", count)

	tmpl := `<div class="bg-green-50 border border-green-200 rounded-md p-4">
		<p class="text-gray-700">Button clicked <strong class="text-green-600">{{.Count}}</strong> times!</p>
		<button class="mt-3 bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
		        hx-post="/api/click" hx-target="#click-counter">Click Me Again!</button>
	</div>`

	t, _ := template.New("click").Parse(tmpl)
	data := struct{ Count int }{Count: count}
	t.Execute(w, data)
}

func HealthcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
