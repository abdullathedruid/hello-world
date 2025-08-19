package main

import (
	"log/slog"
	"net/http"
	"os"

	"hello-world/config"
	"hello-world/routes"
)

// CommitHash is set at build time via ldflags
var CommitHash = "unknown"

func main() {
	// Configure slog with JSON handler
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Load configuration
	cfg := config.Load()

	// Setup routes
	r := routes.SetupRoutes()

	slog.Info("Server starting", "url", "http://localhost:"+cfg.Port, "commit", CommitHash)
	if err := http.ListenAndServe(":"+cfg.Port, r); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
