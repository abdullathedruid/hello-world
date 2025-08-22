package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	// Shared observability components
	observabilityTracer = otel.Tracer("hello-world/middleware")
	observabilityMeter  = otel.Meter("hello-world/middleware")

	// Metrics (initialized once)
	httpRequestDuration      metric.Float64Histogram
	httpRequestSize          metric.Int64Histogram
	httpResponseSize         metric.Int64Histogram
	httpActiveRequests       metric.Int64UpDownCounter
	observabilityInitialized = false
)

// contextKey is an unexported type to avoid collisions in context values
type contextKey int

const (
	requestIDKey contextKey = iota
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

	httpRequestDuration, err = observabilityMeter.Float64Histogram(
		"http.server.request.duration",
		metric.WithDescription("Duration of HTTP server requests"),
		metric.WithUnit("s"),
		metric.WithExplicitBucketBoundaries(0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10),
	)
	if err != nil {
		return err
	}

	httpActiveRequests, err = observabilityMeter.Int64UpDownCounter(
		"http.server.active_requests",
		metric.WithDescription("Number of active HTTP server requests"),
		metric.WithUnit("1"),
	)
	if err != nil {
		return err
	}

	httpRequestSize, err = observabilityMeter.Int64Histogram(
		"http.server.request.size",
		metric.WithDescription("Size of HTTP server requests"),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	httpResponseSize, err = observabilityMeter.Int64Histogram(
		"http.server.response.size",
		metric.WithDescription("Size of HTTP server responses"),
		metric.WithUnit("By"),
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

		// Correlate request ID (accept inbound or generate new)
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestID = generateRequestID()
		}
		ctx := ContextWithRequestID(r.Context(), requestID)

		// Derive route template for low-cardinality span name/attributes
		routeTemplate := r.URL.Path
		if route := mux.CurrentRoute(r); route != nil {
			if tmpl, err := route.GetPathTemplate(); err == nil {
				routeTemplate = tmpl
			}
		}

		// Start tracing span
		ctx, span := observabilityTracer.Start(ctx, r.Method+" "+routeTemplate,
			trace.WithAttributes(
				semconv.HTTPRequestMethodKey.String(r.Method),
				semconv.URLPathKey.String(r.URL.Path),
				semconv.UserAgentOriginalKey.String(r.UserAgent()),
				attribute.String("request_id", requestID),
			),
		)
		defer span.End()

		// Update request with enriched context
		r = r.WithContext(ctx)
		// Return request id for clients
		w.Header().Set("X-Request-Id", requestID)

		// Ensure Content-Length is non-negative
		requestSize := max(r.ContentLength, 0)

		// Increment active requests metric and ensure decrement
		if httpActiveRequests != nil {
			attrs := metric.WithAttributes(
				semconv.HTTPRequestMethodKey.String(r.Method),
				semconv.HTTPRouteKey.String(routeTemplate),
			)
			httpActiveRequests.Add(ctx, 1, attrs)
			defer httpActiveRequests.Add(ctx, -1, attrs)
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
			semconv.HTTPRequestMethodKey.String(r.Method),
			semconv.HTTPRouteKey.String(routeTemplate),
			semconv.HTTPResponseStatusCodeKey.Int(wrapped.statusCode),
		}

		// Update tracing span with response data
		span.SetAttributes(
			append(commonAttrs,
				semconv.HTTPRequestBodySizeKey.Int64(requestSize),
				semconv.HTTPResponseBodySizeKey.Int64(wrapped.responseSize),
			)...,
		)

		// Record metrics
		if observabilityInitialized {
			metricAttrs := metric.WithAttributes(commonAttrs...)

			httpRequestDuration.Record(ctx, duration.Seconds(), metricAttrs)
			httpRequestSize.Record(ctx, requestSize, metricAttrs)
			httpResponseSize.Record(ctx, wrapped.responseSize, metricAttrs)

			// Decrement active requests
			// handled by defer above
		}

	})
}

// ContextWithRequestID adds a request ID to the context
func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// generateRequestID generates a random request ID
func generateRequestID() string {
	bytes := make([]byte, 8)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
