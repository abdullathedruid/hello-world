package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"hello-world/config"
	"hello-world/routes"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	logglobal "go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// CommitHash is set at build time via ldflags
var CommitHash = "unknown"

func main() {
	// Temporary console logger before OTEL init
	tempLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(tempLogger)

	ctx := context.Background()
	shutdown, err := initOtelLogging(ctx)
	if err != nil {
		slog.Warn("OpenTelemetry logging not enabled", "error", err)
	} else {
		slog.Info("OpenTelemetry logging enabled")
		defer func() {
			flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = shutdown(flushCtx)
		}()
	}

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

// initOtelLogging initializes an OTLP HTTP exporter and slog bridge.
func initOtelLogging(ctx context.Context) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		log.Fatalf("OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	exporter, err := otlploghttp.New(
		ctx,
		otlploghttp.WithEndpoint(endpoint),
		otlploghttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String("hello-world"),
			semconv.ServiceVersionKey.String(CommitHash),
		),
	)
	if err != nil {
		return nil, err
	}

	provider := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exporter)),
	)

	logglobal.SetLoggerProvider(provider)

	otelSlog := otelslog.NewLogger(
		"hello-world",
		otelslog.WithLoggerProvider(provider),
	)
	slog.SetDefault(otelSlog)

	return provider.Shutdown, nil
}
