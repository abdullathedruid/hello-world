package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

type PageData struct {
	Title string
	Time  string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/click", clickHandler)
	
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	
	fmt.Println("Server starting on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>HTMX + Go</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link rel="stylesheet" href="/static/style.css">
</head>
<body class="bg-gray-100 min-h-screen">
    <div class="container mx-auto px-4 py-8 max-w-2xl">
        <h1 class="text-4xl font-bold text-center text-gray-800 mb-8">HTMX + Go Demo</h1>
        
        <div class="bg-white rounded-lg shadow-md p-6 mb-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Current Time</h2>
            <div id="time-display">
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
                        hx-get="/time" hx-target="#time-display">Get Current Time</button>
            </div>
        </div>
        
        <div class="bg-white rounded-lg shadow-md p-6">
            <h2 class="text-2xl font-semibold text-gray-700 mb-4">Click Counter</h2>
            <div id="click-counter">
                <button class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
                        hx-post="/click" hx-target="#click-counter">Click Me!</button>
            </div>
        </div>
    </div>
</body>
</html>`
	
	t, _ := template.New("home").Parse(tmpl)
	t.Execute(w, nil)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(w, `<div class="bg-blue-50 border border-blue-200 rounded-md p-4">
		<p class="text-gray-700">Current time: <strong class="text-blue-600">%s</strong></p>
		<button class="mt-3 bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
		        hx-get="/time" hx-target="#time-display">Refresh Time</button>
	</div>`, currentTime)
}

var clickCount = 0

func clickHandler(w http.ResponseWriter, r *http.Request) {
	clickCount++
	fmt.Fprintf(w, `<div class="bg-green-50 border border-green-200 rounded-md p-4">
		<p class="text-gray-700">Button clicked <strong class="text-green-600">%d</strong> times!</p>
		<button class="mt-3 bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded transition duration-200" 
		        hx-post="/click" hx-target="#click-counter">Click Me Again!</button>
	</div>`, clickCount)
}