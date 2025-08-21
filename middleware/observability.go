package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

var (
	// Shared observability components
	observabilityTracer = otel.Tracer("hello-world/middleware")
	observabilityMeter  = otel.Meter("hello-world/middleware")

	// Metrics (initialized once)
	httpRequestsTotal    metric.Int64Counter
	httpRequestDuration  metric.Float64Histogram
	httpRequestSize      metric.Int64Histogram
	httpResponseSize     metric.Int64Histogram
	httpActiveRequests   metric.Int64UpDownCounter
	observabilityInitialized = false
)

// ObservabilityResponseWriter wraps http.ResponseWriter to capture metrics
type ObservabilityResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int64
}

func (w *ObservabilityResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *ObservabilityResponseWriter) Write(b []byte) (int, error) {
	size, err := w.ResponseWriter.Write(b)
	w.responseSize += int64(size)
	return size, err
}

// initObservabilityMetrics initializes metrics once
func initObservabilityMetrics() error {
	if observabilityInitialized {
		return nil
	}

	var err error

	httpRequestsTotal, err = observabilityMeter.Int64Counter(
		"http_server_requests_total",
		metric.WithDescription("Total number of HTTP server requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	httpRequestDuration, err = observabilityMeter.Float64Histogram(
		"http_server_request_duration_seconds",
		metric.WithDescription("Duration of HTTP server requests"),
		metric.WithUnit("s"),
	)
	if err != nil {
		return err
	}

	httpRequestSize, err = observabilityMeter.Int64Histogram(
		"http_server_request_size_bytes",
		metric.WithDescription("Size of HTTP server requests"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	httpResponseSize, err = observabilityMeter.Int64Histogram(
		"http_server_response_size_bytes",
		metric.WithDescription("Size of HTTP server responses"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	httpActiveRequests, err = observabilityMeter.Int64UpDownCounter(
		"http_server_active_requests",
		metric.WithDescription("Number of active HTTP server requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	observabilityInitialized = true
	return nil
}

// ObservabilityMiddleware provides unified logging, tracing, and metrics
func ObservabilityMiddleware(next http.Handler) http.Handler {
	// Initialize metrics on first use
	if err := initObservabilityMetrics(); err != nil {
		slog.Warn("Failed to initialize observability metrics", "error", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Generate request ID for correlation
		requestID := generateRequestID()
		ctx := ContextWithRequestID(r.Context(), requestID)

		// Start tracing span
		ctx, span := observabilityTracer.Start(ctx, "http_request",
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.Path),
				attribute.String("http.user_agent", r.UserAgent()),
				attribute.String("request_id", requestID),
			),
		)
		defer span.End()

		// Update request with enriched context
		r = r.WithContext(ctx)

		// Calculate request size from Content-Length header
		requestSize := r.ContentLength
		if requestSize < 0 {
			requestSize = 0 // Unknown content length
		}

		// Increment active requests metric
		if httpActiveRequests != nil {
			httpActiveRequests.Add(ctx, 1,
				metric.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.route", r.URL.Path),
				),
			)
		}

		// Wrap response writer for metrics collection
		wrapped := &ObservabilityResponseWriter{
			ResponseWriter: w,
			statusCode:     200, // Default status code
		}

		// Process request
		next.ServeHTTP(wrapped, r)

		// Calculate duration
		duration := time.Since(start)

		// Common attributes for all observability signals (following OpenTelemetry semantic conventions)
		commonAttrs := []attribute.KeyValue{
			attribute.String("http.method", r.Method),
			attribute.String("http.route", r.URL.Path),
			attribute.Int("http.status_code", wrapped.statusCode),
		}

		// Update tracing span with response data
		span.SetAttributes(
			append(commonAttrs,
				attribute.Int64("http.duration_ms", duration.Milliseconds()),
				attribute.Int64("http.request_size", requestSize),
				attribute.Int64("http.response_size", wrapped.responseSize),
			)...,
		)

		// Record metrics
		if observabilityInitialized {
			metricAttrs := metric.WithAttributes(commonAttrs...)
			
			httpRequestsTotal.Add(ctx, 1, metricAttrs)
			httpRequestDuration.Record(ctx, duration.Seconds(), metricAttrs)
			httpRequestSize.Record(ctx, requestSize, metricAttrs)
			httpResponseSize.Record(ctx, wrapped.responseSize, metricAttrs)
			
			// Decrement active requests
			httpActiveRequests.Add(ctx, -1,
				metric.WithAttributes(
					attribute.String("http.method", r.Method),
					attribute.String("http.route", r.URL.Path),
				),
			)
		}

		// Structured logging with trace correlation
		slog.InfoContext(ctx, "Request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"status_code", wrapped.statusCode,
			"duration", duration,
			"response_size", wrapped.responseSize,
		)
	})
}

// ContextWithRequestID adds a request ID to the context
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, "request_id", requestID)
}

// generateRequestID generates a random request ID
func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}