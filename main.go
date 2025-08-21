package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"hello-world/config"
	"hello-world/middleware"
	"hello-world/routes"

	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	logglobal "go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// CommitHash is set at build time via ldflags
var CommitHash = "unknown"

func main() {
	// Temporary console logger before OTEL init
	tempLogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(tempLogger)

	ctx := context.Background()

	// Initialize tracing
	traceShutdown, err := initOtelTracing(ctx)
	if err != nil {
		slog.Warn("OpenTelemetry tracing not enabled", "error", err)
	} else {
		slog.Info("OpenTelemetry tracing enabled")
		defer func() {
			flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = traceShutdown(flushCtx)
		}()
	}

	// Initialize metrics
	metricsShutdown, err := initOtelMetrics(ctx)
	if err != nil {
		slog.Warn("OpenTelemetry metrics not enabled", "error", err)
	} else {
		slog.Info("OpenTelemetry metrics enabled")
		defer func() {
			flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = metricsShutdown(flushCtx)
		}()
	}

	// Initialize logging with trace integration
	logShutdown, err := initOtelLogging(ctx)
	if err != nil {
		slog.Warn("OpenTelemetry logging not enabled", "error", err)
	} else {
		slog.Info("OpenTelemetry logging enabled")
		defer func() {
			flushCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = logShutdown(flushCtx)
		}()
	}

	// Load configuration
	cfg := config.Load()

	// Setup routes and HTTP server
	r := routes.SetupRoutes()
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		slog.Info("Server starting", "url", "http://localhost:"+cfg.Port, "commit", CommitHash)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
		}
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	slog.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("Graceful shutdown failed", "error", err)
	}
}

// initOtelLogging initializes an OTLP HTTP exporter and slog bridge.
func initOtelLogging(ctx context.Context) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	// TLS/insecure config via env
	insecure := otlpInsecure()
	logOpts := []otlploghttp.Option{otlploghttp.WithEndpoint(endpoint)}
	if insecure {
		logOpts = append(logOpts, otlploghttp.WithInsecure())
	} else {
		logOpts = append(logOpts, otlploghttp.WithTLSClientConfig(&tls.Config{MinVersion: tls.VersionTLS12}))
	}

	exporter, err := otlploghttp.New(ctx, logOpts...)
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

	// Wrap with TraceHandler to add trace IDs to logs
	traceHandler := middleware.NewTraceHandler(otelSlog.Handler())
	tracedLogger := slog.New(traceHandler)
	slog.SetDefault(tracedLogger)

	return provider.Shutdown, nil
}

// initOtelTracing initializes OpenTelemetry tracing
func initOtelTracing(ctx context.Context) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	// Create trace exporter
	insecure := otlpInsecure()
	traceOpts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(endpoint)}
	if insecure {
		traceOpts = append(traceOpts, otlptracehttp.WithInsecure())
	} else {
		traceOpts = append(traceOpts, otlptracehttp.WithTLSClientConfig(&tls.Config{MinVersion: tls.VersionTLS12}))
	}
	traceExporter, err := otlptracehttp.New(ctx, traceOpts...)
	if err != nil {
		return nil, err
	}

	// Create resource
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

	// Create trace provider
	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(envOrDefaultSampler()),
	)

	// Set global trace provider and composite propagator
	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return traceProvider.Shutdown, nil
}

// initOtelMetrics initializes OpenTelemetry metrics
func initOtelMetrics(ctx context.Context) (func(context.Context) error, error) {
	endpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if endpoint == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	// Create metrics exporter
	insecure := otlpInsecure()
	metricOpts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(endpoint),
		otlpmetrichttp.WithURLPath("/v1/metrics"),
	}
	if insecure {
		metricOpts = append(metricOpts, otlpmetrichttp.WithInsecure())
	} else {
		metricOpts = append(metricOpts, otlpmetrichttp.WithTLSClientConfig(&tls.Config{MinVersion: tls.VersionTLS12}))
	}
	metricsExporter, err := otlpmetrichttp.New(ctx, metricOpts...)
	if err != nil {
		return nil, err
	}

	// Create resource
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

	// Create metrics provider
	reader := sdkmetric.NewPeriodicReader(
		metricsExporter,
		sdkmetric.WithInterval(30*time.Second), // Export every 30 seconds
	)

	metricsProvider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(reader),
		sdkmetric.WithResource(res),
	)

	// Set global metrics provider
	otel.SetMeterProvider(metricsProvider)

	return metricsProvider.Shutdown, nil
}

// otlpInsecure returns whether to use insecure transport to the OTLP endpoint. Defaults to true.
func otlpInsecure() bool {
	v := os.Getenv("OTEL_EXPORTER_OTLP_INSECURE")
	if v == "" {
		return true
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return true
	}
	return b
}

// envOrDefaultSampler selects a sampler from env vars or defaults to ParentBased(10%).
func envOrDefaultSampler() sdktrace.Sampler {
	if arg := os.Getenv("OTEL_TRACES_SAMPLER_ARG"); arg != "" {
		if ratio, err := strconv.ParseFloat(arg, 64); err == nil {
			if ratio < 0 {
				ratio = 0
			}
			if ratio > 1 {
				ratio = 1
			}
			return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))
		}
	}
	return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.10))
}
