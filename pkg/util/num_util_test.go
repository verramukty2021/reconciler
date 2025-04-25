package util

import (
	"testing"
)

func TestRoundToTwoDecimals(t *testing.T) {
	tests := []struct {
		input    float64
		expected float64
	}{
		{123.456, 123.46},
		{123.454, 123.45},
		{0.005, 0.01},
		{-1.2345, -1.23},
		{-1.2356, -1.24},
	}

	for _, tt := range tests {
		result := RoundToTwoDecimals(tt.input)
		if result != tt.expected {
			t.Errorf("RoundToTwoDecimals(%v) = %v; want %v", tt.input, result, tt.expected)
		}
	}
}
