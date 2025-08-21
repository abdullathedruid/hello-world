package middleware

import (
	"context"
	"log/slog"

	"go.opentelemetry.io/otel/trace"
)

// TraceHandler wraps an slog.Handler to automatically inject trace information
type TraceHandler struct {
	slog.Handler
}

// NewTraceHandler creates a new TraceHandler that wraps the provided handler
func NewTraceHandler(handler slog.Handler) *TraceHandler {
	return &TraceHandler{Handler: handler}
}

// Handle processes log records and automatically adds trace information from context
func (h *TraceHandler) Handle(ctx context.Context, r slog.Record) error {
	// Automatically extract trace data from context
	if span := trace.SpanFromContext(ctx); span.SpanContext().IsValid() {
		r.AddAttrs(
			slog.String("trace_id", span.SpanContext().TraceID().String()),
			slog.String("span_id", span.SpanContext().SpanID().String()),
		)
	}

	// Extract request ID if present
	if reqID := getRequestID(ctx); reqID != "" {
		r.AddAttrs(slog.String("request_id", reqID))
	}

	return h.Handler.Handle(ctx, r)
}

// getRequestID extracts request ID from context
func getRequestID(ctx context.Context) string {
	if reqID := ctx.Value("request_id"); reqID != nil {
		if id, ok := reqID.(string); ok {
			return id
		}
	}
	return ""
}
