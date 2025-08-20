package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"strings"
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
// Env vars used:
// - OTEL_EXPORTER_OTLP_ENDPOINT (default: http://localhost:4318)
// - OTEL_EXPORTER_OTLP_HEADERS  (key1=value1,key2=value2)
func initOtelLogging(ctx context.Context) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:4318"
	}

	// Normalize endpoint for otlploghttp: use host:port, stripping scheme and path if present
	if strings.HasPrefix(endpoint, "http://") || strings.HasPrefix(endpoint, "https://") {
		noScheme := strings.TrimPrefix(strings.TrimPrefix(endpoint, "http://"), "https://")
		if slash := strings.Index(noScheme, "/"); slash != -1 {
			endpoint = noScheme[:slash]
		} else {
			endpoint = noScheme
		}
	}

	headersEnv := os.Getenv("OTEL_EXPORTER_OTLP_HEADERS")
	headers := map[string]string{}
	if headersEnv != "" {
		for _, pair := range strings.Split(headersEnv, ",") {
			kv := strings.SplitN(strings.TrimSpace(pair), "=", 2)
			if len(kv) == 2 {
				headers[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
			}
		}
	}

	// If a CLICKSTACK_API_KEY is present, set Authorization header unless already provided
	if apiKey := os.Getenv("CLICKSTACK_API_KEY"); apiKey != "" {
		if _, ok := headers["authorization"]; !ok {
			if _, ok2 := headers["Authorization"]; !ok2 {
				headers["authorization"] = apiKey
			}
		}
	}

	exporter, err := otlploghttp.New(
		ctx,
		otlploghttp.WithEndpoint(endpoint),
		otlploghttp.WithInsecure(),
		otlploghttp.WithHeaders(headers),
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
