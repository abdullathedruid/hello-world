package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"hello-world/handlers"
)

func BenchmarkHomeHandler(b *testing.B) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.HomeHandler)
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkTimeHandler(b *testing.B) {
	req, err := http.NewRequest("GET", "/time", nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.TimeFragmentHandler)
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkClickHandler(b *testing.B) {
	req, err := http.NewRequest("POST", "/click", nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.ClickFragmentHandler)
		handler.ServeHTTP(rr, req)
	}
}

func BenchmarkDebugHandler(b *testing.B) {
	req, err := http.NewRequest("GET", "/debug", nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(handlers.DebugHandler)
		handler.ServeHTTP(rr, req)
	}
}
