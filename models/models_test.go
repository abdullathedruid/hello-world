package models

import (
	"testing"
)

func TestPageDataStruct(t *testing.T) {
	data := PageData{
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

func TestClickDataStruct(t *testing.T) {
	data := ClickData{
		Count: 42,
	}

	if data.Count != 42 {
		t.Errorf("Expected Count to be 42, got %d", data.Count)
	}
}

func TestPageDataEmpty(t *testing.T) {
	var data PageData

	if data.Title != "" {
		t.Errorf("Expected empty Title, got %s", data.Title)
	}

	if data.Time != "" {
		t.Errorf("Expected empty Time, got %s", data.Time)
	}
}

func TestClickDataEmpty(t *testing.T) {
	var data ClickData

	if data.Count != 0 {
		t.Errorf("Expected Count to be 0, got %d", data.Count)
	}
}