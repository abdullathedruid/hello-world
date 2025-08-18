package handlers

import (
	"html/template"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
	"hello-world/services"
)

var clickService = services.NewClickService()

func TimeFragmentHandler(w http.ResponseWriter, r *http.Request) {
	logrus.Info("Time endpoint accessed")
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
	logrus.WithField("count", count).Info("Button clicked")
	
	tmpl := `<div class="bg-green-50 border border-green-200 rounded-md p-4">
		<p class="text-gray-700">Button clicked <strong class="text-green-600">{{.Count}}</strong> times!</p>
		<button class="mt-3 bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
		        hx-post="/api/click" hx-target="#click-counter">Click Me Again!</button>
	</div>`
	
	t, _ := template.New("click").Parse(tmpl)
	data := struct{ Count int }{Count: count}
	t.Execute(w, data)
}