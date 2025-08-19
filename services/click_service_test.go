package services

import (
	"testing"
)

func TestNewClickService(t *testing.T) {
	service := NewClickService()
	if service == nil {
		t.Error("NewClickService should return a valid service")
	}
}

func TestIncrementClick(t *testing.T) {
	service := NewClickService()

	// Reset to ensure clean state
	service.Reset()

	// First increment
	count := service.IncrementClick()
	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	// Second increment
	count = service.IncrementClick()
	if count != 2 {
		t.Errorf("Expected count 2, got %d", count)
	}

	// Third increment
	count = service.IncrementClick()
	if count != 3 {
		t.Errorf("Expected count 3, got %d", count)
	}
}

func TestGetCount(t *testing.T) {
	service := NewClickService()

	// Reset to ensure clean state
	service.Reset()

	// Initial count should be 0
	count := service.GetCount()
	if count != 0 {
		t.Errorf("Expected initial count 0, got %d", count)
	}

	// Increment and check
	service.IncrementClick()
	count = service.GetCount()
	if count != 1 {
		t.Errorf("Expected count 1 after increment, got %d", count)
	}
}

func TestReset(t *testing.T) {
	service := NewClickService()

	// Reset to ensure clean state
	service.Reset()

	// Increment a few times
	service.IncrementClick()
	service.IncrementClick()
	service.IncrementClick()

	// Verify count is not 0
	count := service.GetCount()
	if count != 3 {
		t.Errorf("Expected count 3 before reset, got %d", count)
	}

	// Reset and verify
	service.Reset()
	count = service.GetCount()
	if count != 0 {
		t.Errorf("Expected count 0 after reset, got %d", count)
	}
}

func TestMultipleServices(t *testing.T) {
	// Note: Since clickCount is a global variable, multiple services share state
	// This tests the current implementation behavior
	service1 := NewClickService()
	service2 := NewClickService()

	service1.Reset()

	// Increment with first service
	count1 := service1.IncrementClick()
	if count1 != 1 {
		t.Errorf("Expected count 1 from service1, got %d", count1)
	}

	// Check with second service (should see same count due to shared state)
	count2 := service2.GetCount()
	if count2 != 1 {
		t.Errorf("Expected count 1 from service2, got %d", count2)
	}

	// Increment with second service
	count2 = service2.IncrementClick()
	if count2 != 2 {
		t.Errorf("Expected count 2 from service2, got %d", count2)
	}
}
