package util

import (
	"testing"
)

func TestDeleteFromSliceMap(t *testing.T) {
	// Helper to build test data
	buildSlice := func(items ...string) []interface{} {
		s := make([]interface{}, len(items))
		for i, v := range items {
			s[i] = v
		}
		return s
	}

	t.Run("Valid deletion", func(t *testing.T) {
		mapData := map[string][]interface{}{
			"testKey": buildSlice("A", "B", "C", "D", "E"),
		}
		ok, err := DeleteFromSliceMap(mapData, "testKey", 1, 4)
		if !ok || err != nil {
			t.Errorf("Expected deletion to succeed, got ok=%v, err=%v", ok, err)
		}
		expected := buildSlice("A", "E")
		got := mapData["testKey"]
		if len(got) != len(expected) {
			t.Errorf("Expected length %d, got %d", len(expected), len(got))
		}
		for i := range expected {
			if got[i] != expected[i] {
				t.Errorf("Expected value at index %d to be %v, got %v", i, expected[i], got[i])
			}
		}
	})

	t.Run("Invalid index range", func(t *testing.T) {
		mapData := map[string][]interface{}{
			"testKey": buildSlice("X", "Y", "Z"),
		}
		ok, err := DeleteFromSliceMap(mapData, "testKey", 2, 1)
		if ok {
			t.Errorf("Expected deletion to fail due to invalid index range")
		}
		if err == nil {
			t.Error("Expected error for invalid indices, got nil")
		}
	})

	t.Run("Key not found", func(t *testing.T) {
		mapData := map[string][]interface{}{
			"otherKey": buildSlice("1", "2"),
		}
		ok, err := DeleteFromSliceMap(mapData, "missingKey", 0, 1)
		if ok {
			t.Error("Expected deletion to return false when key is not found")
		}
		if err != nil {
			t.Errorf("Expected no error when key is missing, got %v", err)
		}
	})

	t.Run("Boundary index test", func(t *testing.T) {
		mapData := map[string][]interface{}{
			"testKey": buildSlice("A", "B", "C"),
		}
		ok, err := DeleteFromSliceMap(mapData, "testKey", 0, 3)
		if !ok || err != nil {
			t.Errorf("Expected full slice deletion, got ok=%v, err=%v", ok, err)
		}
		if len(mapData["testKey"]) != 0 {
			t.Errorf("Expected empty slice, got %v", mapData["testKey"])
		}
	})
}
