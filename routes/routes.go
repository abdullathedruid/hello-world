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

	// Apply unified observability middleware to a subrouter for all other routes
	observed := r.NewRoute().Subrouter()
	observed.Use(middleware.ObservabilityMiddleware)

	// Full page routes
	observed.HandleFunc("/", handlers.HomeHandler).Methods("GET")
	observed.HandleFunc("/debug", handlers.DebugHandler).Methods("GET")

	// API routes for HTMX fragments
	api := observed.PathPrefix("/api").Subrouter()
	api.HandleFunc("/time", handlers.TimeFragmentHandler).Methods("GET")
	api.HandleFunc("/click", handlers.ClickFragmentHandler).Methods("POST")

	// Static files
	observed.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	return r
}
