package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestTimeFragmentHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/time", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(TimeFragmentHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Current time:") {
		t.Errorf("handler should return current time fragment")
	}

	if !strings.Contains(body, "Refresh Time") {
		t.Errorf("handler should contain refresh button")
	}

	if !strings.Contains(body, "bg-blue-50") {
		t.Errorf("handler should contain proper CSS classes")
	}
}

func TestClickFragmentHandler(t *testing.T) {
	// Reset click service for isolated test
	ResetClickService()

	req, err := http.NewRequest("POST", "/api/click", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ClickFragmentHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Button clicked") {
		t.Errorf("handler should show button clicked message")
	}

	if !strings.Contains(body, "1") {
		t.Errorf("handler should show click count of 1")
	}

	if !strings.Contains(body, "Click Me Again!") {
		t.Errorf("handler should contain click button")
	}

	if !strings.Contains(body, "bg-green-50") {
		t.Errorf("handler should contain proper CSS classes")
	}
}

func TestClickFragmentHandlerMultipleClicks(t *testing.T) {
	// Reset click service for isolated test
	ResetClickService()

	// First click
	req1, err := http.NewRequest("POST", "/api/click", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	handler := http.HandlerFunc(ClickFragmentHandler)
	handler.ServeHTTP(rr1, req1)

	// Second click
	req2, err := http.NewRequest("POST", "/api/click", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)

	body2 := rr2.Body.String()
	if !strings.Contains(body2, "2") {
		t.Errorf("handler should show click count of 2 after second click")
	}
}

func TestResetClickService(t *testing.T) {
	// Increment click count
	req, err := http.NewRequest("POST", "/api/click", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ClickFragmentHandler)
	handler.ServeHTTP(rr, req)

	// Reset
	ResetClickService()

	// Check that next click shows count 1
	req2, err := http.NewRequest("POST", "/api/click", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr2 := httptest.NewRecorder()
	handler.ServeHTTP(rr2, req2)

	body := rr2.Body.String()
	if !strings.Contains(body, "1") {
		t.Errorf("handler should show click count of 1 after reset")
	}
}
