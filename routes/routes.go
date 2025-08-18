package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"hello-world/handlers"
	"hello-world/middleware"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	
	// Apply middleware
	r.Use(middleware.LoggingMiddleware)
	
	// Full page routes
	r.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	r.HandleFunc("/debug", handlers.DebugHandler).Methods("GET")
	
	// API routes for HTMX fragments
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/time", handlers.TimeFragmentHandler).Methods("GET")
	api.HandleFunc("/click", handlers.ClickFragmentHandler).Methods("POST")
	
	// Static files
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	
	return r
}