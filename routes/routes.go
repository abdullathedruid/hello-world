package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"hello-world/handlers"
	"hello-world/middleware"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Healthcheck endpoint without middleware
	r.HandleFunc("/health", handlers.HealthcheckHandler).Methods("GET")

	// Apply middleware to a subrouter for all other routes
	logged := r.NewRoute().Subrouter()
	logged.Use(middleware.LoggingMiddleware)

	// Full page routes
	logged.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	logged.HandleFunc("/debug", handlers.DebugHandler).Methods("GET")

	// API routes for HTMX fragments
	api := logged.PathPrefix("/api").Subrouter()
	api.HandleFunc("/time", handlers.TimeFragmentHandler).Methods("GET")
	api.HandleFunc("/click", handlers.ClickFragmentHandler).Methods("POST")

	// Static files
	logged.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return r
}
