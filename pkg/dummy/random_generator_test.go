package dummy

import (
	"testing"
	"time"
)

// Test GenerateRandomID
func TestGenerateRandomID(t *testing.T) {
	id := GenerateRandomID()
	if len(id) != idLength {
		t.Errorf("Expected ID length %d, got %d", idLength, len(id))
	}
	for _, c := range id {
		if !(('A' <= c && c <= 'Z') || ('0' <= c && c <= '9')) {
			t.Errorf("Invalid character in ID: %c", c)
		}
	}
}

// Test GenerateRandomAmount
func TestGenerateRandomAmount(t *testing.T) {
	min := 10.0
	max := 100.0

	for i := 0; i < 100; i++ {
		amt := GenerateRandomAmount(min, max)
		if amt < min || amt > max {
			t.Errorf("Amount %.2f out of bounds [%f, %f]", amt, min, max)
		}
	}
}

// Test GenerateRandomTrxType
func TestGenerateRandomTrxType(t *testing.T) {
	for i := 0; i < 100; i++ {
		tt := GenerateRandomTrxType()
		if tt != "DEBIT" && tt != "CREDIT" {
			t.Errorf("Unexpected transaction type: %s", tt)
		}
	}
}

// Test GenerateRandomTime
func TestGenerateRandomTime(t *testing.T) {
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	for i := 0; i < 100; i++ {
		rt := GenerateRandomTime(start, end)
		if rt.Before(start) || rt.After(end) {
			t.Errorf("Generated time %v out of range [%v, %v]", rt, start, end)
		}
	}
}
