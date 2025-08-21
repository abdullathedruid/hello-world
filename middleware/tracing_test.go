package middleware

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"go.opentelemetry.io/otel"
)

func TestTraceHandler(t *testing.T) {
	// Create a test handler that captures log records
	var lastRecord slog.Record
	testHandler := &testLogHandler{
		handleFunc: func(ctx context.Context, r slog.Record) error {
			lastRecord = r
			return nil
		},
	}

	// Create trace handler
	traceHandler := NewTraceHandler(testHandler)

	// Create a test context with a span
	tracer := otel.Tracer("test")
	ctx, span := tracer.Start(context.Background(), "test-span")
	defer span.End()

	// Add request ID to context
	ctx = ContextWithRequestID(ctx, "test-request-123")

	// Create a log record
	record := slog.NewRecord(time.Now(), slog.LevelInfo, "test message", 0)

	// Handle the record
	err := traceHandler.Handle(ctx, record)
	if err != nil {
		t.Fatalf("Handle() error = %v", err)
	}

	// Check if trace information was added
	found := false
	lastRecord.Attrs(func(a slog.Attr) bool {
		if a.Key == "trace_id" || a.Key == "span_id" || a.Key == "request_id" {
			found = true
		}
		return true
	})

	if !found {
		t.Error("Expected trace information to be added to log record")
	}
}

// testLogHandler is a test implementation of slog.Handler
type testLogHandler struct {
	handleFunc func(context.Context, slog.Record) error
}

func (h *testLogHandler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *testLogHandler) Handle(ctx context.Context, r slog.Record) error {
	if h.handleFunc != nil {
		return h.handleFunc(ctx, r)
	}
	return nil
}

func (h *testLogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *testLogHandler) WithGroup(name string) slog.Handler {
	return h
}

func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()

	if id1 == id2 {
		t.Error("Expected different request IDs")
	}

	if len(id1) != 16 { // 8 bytes = 16 hex characters
		t.Errorf("Expected request ID length 16, got %d", len(id1))
	}
}

func TestContextWithRequestID(t *testing.T) {
	ctx := context.Background()
	requestID := "test-123"

	ctx = ContextWithRequestID(ctx, requestID)

	if got := getRequestID(ctx); got != requestID {
		t.Errorf("Expected request ID %s, got %s", requestID, got)
	}
}
