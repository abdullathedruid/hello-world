package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Test default port when PORT env var is not set
	os.Unsetenv("PORT")
	config := Load()
	if config.Port != "8080" {
		t.Errorf("Expected default port 8080, got %s", config.Port)
	}

	// Test custom port when PORT env var is set
	os.Setenv("PORT", "9000")
	config = Load()
	if config.Port != "9000" {
		t.Errorf("Expected port 9000, got %s", config.Port)
	}

	// Clean up
	os.Unsetenv("PORT")
}

func TestConfigStruct(t *testing.T) {
	config := &Config{Port: "3000"}
	if config.Port != "3000" {
		t.Errorf("Expected port 3000, got %s", config.Port)
	}
}
