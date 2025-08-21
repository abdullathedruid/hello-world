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
	"go.opentelemetry.io/otel/trace"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	tracer := otel.Tracer("hello-world/middleware")
	
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Start a new span for this request
		ctx, span := tracer.Start(r.Context(), "http_request",
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.url", r.URL.Path),
				attribute.String("http.user_agent", r.UserAgent()),
			),
		)
		defer span.End()

		// Add request ID to context for tracing
		ctx = ContextWithRequestID(ctx, generateRequestID())
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

		duration := time.Since(start)
		
		// Add span attributes for response
		span.SetAttributes(
			attribute.Int64("http.duration_ms", duration.Milliseconds()),
		)

		// Log with context that includes trace information
		slog.InfoContext(ctx, "Request processed",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration)
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
