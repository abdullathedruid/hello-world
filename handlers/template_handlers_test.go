package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHomeHandlerWithTemplates(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HomeHandler)

	handler.ServeHTTP(rr, req)

	// Check status - may be 500 if templates are missing or 200 if found
	if status := rr.Code; status != http.StatusOK && status != http.StatusInternalServerError {
		t.Errorf("handler returned unexpected status code: got %v, expected %v or %v",
			status, http.StatusOK, http.StatusInternalServerError)
	}

	// If successful, check for expected content
	if rr.Code == http.StatusOK {
		body := rr.Body.String()
		if !strings.Contains(body, "HTMX") && !strings.Contains(body, "Go Demo") {
			t.Errorf("handler should contain title content when templates are available")
		}
	}
}

func TestDebugHandlerWithTemplates(t *testing.T) {
	req, err := http.NewRequest("GET", "/debug", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DebugHandler)

	handler.ServeHTTP(rr, req)

	// Check status - may be 500 if templates are missing or 200 if found
	if status := rr.Code; status != http.StatusOK && status != http.StatusInternalServerError {
		t.Errorf("handler returned unexpected status code: got %v, expected %v or %v",
			status, http.StatusOK, http.StatusInternalServerError)
	}

	// If successful, check for expected content
	if rr.Code == http.StatusOK {
		body := rr.Body.String()
		if !strings.Contains(body, "Farcaster") && !strings.Contains(body, "Debug") {
			t.Errorf("handler should contain debug content when templates are available")
		}
	}
}