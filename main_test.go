package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hello-world/handlers"
	"hello-world/models"
	"hello-world/routes"
)

func TestHomeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.HomeHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "HTMX") || !strings.Contains(body, "Go Demo") {
		t.Errorf("handler returned unexpected body: should contain HTMX and Go Demo")
	}

	if !strings.Contains(body, "Current Time") {
		t.Errorf("handler should contain 'Current Time' section")
	}

	if !strings.Contains(body, "Click Counter") {
		t.Errorf("handler should contain 'Click Counter' section")
	}
}

func TestTimeHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/time", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.TimeFragmentHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Current time:") {
		t.Errorf("handler should return current time")
	}

	if !strings.Contains(body, "Refresh Time") {
		t.Errorf("handler should contain refresh button")
	}
}

func TestClickHandler(t *testing.T) {
	// Reset click service for isolated test
	handlers.ResetClickService()

	req, err := http.NewRequest("POST", "/click", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.ClickFragmentHandler)

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
}

func TestDebugHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/debug", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.DebugHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "Farcaster MiniApp Debug") {
		t.Errorf("handler should contain debug page title")
	}

	if !strings.Contains(body, "Context Information") {
		t.Errorf("handler should contain context information section")
	}

	if !strings.Contains(body, "SDK Status") {
		t.Errorf("handler should contain SDK status section")
	}
}

func TestRoutes(t *testing.T) {
	r := routes.SetupRoutes()

	tests := []struct {
		method         string
		path           string
		expectedStatus int
	}{
		{"GET", "/", http.StatusOK},
		{"GET", "/debug", http.StatusOK},
		{"GET", "/api/time", http.StatusOK},
		{"POST", "/api/click", http.StatusOK},
		{"GET", "/nonexistent", http.StatusNotFound},
		{"POST", "/", http.StatusMethodNotAllowed},
		{"GET", "/api/click", http.StatusMethodNotAllowed},
	}

	for _, test := range tests {
		req, err := http.NewRequest(test.method, test.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code != test.expectedStatus {
			t.Errorf("Route %s %s: expected status %d, got %d",
				test.method, test.path, test.expectedStatus, rr.Code)
		}
	}
}

func TestPageDataStruct(t *testing.T) {
	data := models.PageData{
		Title: "Test Title",
		Time:  "2023-01-01 12:00:00",
	}

	if data.Title != "Test Title" {
		t.Errorf("Expected Title to be 'Test Title', got %s", data.Title)
	}

	if data.Time != "2023-01-01 12:00:00" {
		t.Errorf("Expected Time to be '2023-01-01 12:00:00', got %s", data.Time)
	}
}
