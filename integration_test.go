package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

func TestFullServerIntegration(t *testing.T) {
	// Create router with all routes
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler).Methods("GET")
	r.HandleFunc("/debug", debugHandler).Methods("GET")
	r.HandleFunc("/time", timeHandler).Methods("GET")
	r.HandleFunc("/click", clickHandler).Methods("POST")

	// Create test server
	server := httptest.NewServer(r)
	defer server.Close()

	// Test home page
	resp, err := http.Get(server.URL + "/")
	if err != nil {
		t.Fatalf("Failed to get home page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	if !strings.Contains(string(body), "HTMX + Go Demo") {
		t.Error("Home page should contain title")
	}

	// Test debug endpoint
	resp, err = http.Get(server.URL + "/debug")
	if err != nil {
		t.Fatalf("Failed to get debug page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for debug page, got %d", resp.StatusCode)
	}

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read debug response: %v", err)
	}

	if !strings.Contains(string(body), "Farcaster MiniApp Debug") {
		t.Error("Debug page should contain debug title")
	}

	// Test time endpoint
	resp, err = http.Get(server.URL + "/time")
	if err != nil {
		t.Fatalf("Failed to get time: %v", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read time response: %v", err)
	}

	if !strings.Contains(string(body), "Current time:") {
		t.Error("Time endpoint should return current time")
	}

	// Test click endpoint
	resp, err = http.Post(server.URL+"/click", "", nil)
	if err != nil {
		t.Fatalf("Failed to post to click: %v", err)
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read click response: %v", err)
	}

	if !strings.Contains(string(body), "Button clicked") {
		t.Error("Click endpoint should show button clicked message")
	}
}

func TestConcurrentRequests(t *testing.T) {
	r := mux.NewRouter()
	r.HandleFunc("/time", timeHandler).Methods("GET")

	server := httptest.NewServer(r)
	defer server.Close()

	// Reset click count for concurrent test
	originalCount := clickCount
	clickCount = 0
	defer func() { clickCount = originalCount }()

	// Test concurrent requests
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			resp, err := http.Get(server.URL + "/time")
			if err != nil {
				t.Errorf("Concurrent request failed: %v", err)
			} else {
				resp.Body.Close()
				if resp.StatusCode != http.StatusOK {
					t.Errorf("Expected status 200, got %d", resp.StatusCode)
				}
			}
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < 10; i++ {
		select {
		case <-done:
			// Request completed
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout waiting for concurrent requests")
		}
	}
}
