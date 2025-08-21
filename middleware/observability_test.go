package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestObservabilityMiddleware(t *testing.T) {
	// Create a test handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	})

	// Wrap with observability middleware
	observabilityHandler := ObservabilityMiddleware(handler)

	// Create test request
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute request
	observabilityHandler.ServeHTTP(w, req)

	// Verify response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	if w.Body.String() != "test response" {
		t.Errorf("Expected 'test response', got %s", w.Body.String())
	}
}

func TestObservabilityResponseWriter(t *testing.T) {
	// Test our response writer wrapper
	w := httptest.NewRecorder()
	rw := &ObservabilityResponseWriter{
		ResponseWriter: w,
		statusCode:     200,
	}

	// Test WriteHeader
	rw.WriteHeader(404)
	if rw.statusCode != 404 {
		t.Errorf("Expected status code 404, got %d", rw.statusCode)
	}

	// Test Write
	data := []byte("test data")
	n, err := rw.Write(data)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if n != len(data) {
		t.Errorf("Expected %d bytes written, got %d", len(data), n)
	}
	if rw.responseSize != int64(len(data)) {
		t.Errorf("Expected response size %d, got %d", len(data), rw.responseSize)
	}
}

func TestContextWithRequestID(t *testing.T) {
	// Create a context with request ID
	requestID := "test-request-id"
	req := httptest.NewRequest("GET", "/", nil)
	ctx := ContextWithRequestID(req.Context(), requestID)

	// Verify request ID can be retrieved
	retrievedID := ctx.Value("request_id")
	if retrievedID != requestID {
		t.Errorf("Expected request ID '%s', got '%v'", requestID, retrievedID)
	}
}

func TestGenerateRequestID(t *testing.T) {
	// Generate multiple request IDs
	id1 := generateRequestID()
	id2 := generateRequestID()

	// Verify they are different
	if id1 == id2 {
		t.Error("Request IDs should be unique")
	}

	// Verify they are not empty
	if id1 == "" || id2 == "" {
		t.Error("Request IDs should not be empty")
	}

	// Verify they are hex strings of expected length
	expectedLength := 16 // 8 bytes = 16 hex characters
	if len(id1) != expectedLength || len(id2) != expectedLength {
		t.Errorf("Expected request ID length %d, got %d and %d", expectedLength, len(id1), len(id2))
	}
}
