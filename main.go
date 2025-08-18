package main

import (
	"net/http"

	"github.com/sirupsen/logrus"
	"hello-world/config"
	"hello-world/routes"
)

func main() {
	// Configure logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)

	// Load configuration
	cfg := config.Load()

	// Setup routes
	r := routes.SetupRoutes()

	logrus.Info("Server starting on http://localhost:" + cfg.Port)
	logrus.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
