package main

import (
	"testing"
)

// Helper to generate test data
func makeTrxMap(key string, count int, value interface{}) map[string][]interface{} {
	m := make(map[string][]interface{})
	rows := make([]interface{}, count)
	for i := 0; i < count; i++ {
		rows[i] = value
	}
	m[key] = rows
	return m
}

func TestDeleteMatchedTrx_EqualCounts(t *testing.T) {
	systemTrx := makeTrxMap("trx123", 2, "sysRow")
	bankTrx := makeTrxMap("trx123", 2, "bankRow")

	count := deleteMatchedTrx(systemTrx, bankTrx)

	if count != 4 {
		t.Errorf("Expected 4 matched transactions, got %d", count)
	}
	if len(systemTrx) != 0 || len(bankTrx) != 0 {
		t.Errorf("Expected both maps to be empty after match")
	}
}

func TestDeleteMatchedTrx_MoreBankTrx(t *testing.T) {
	systemTrx := makeTrxMap("trx123", 2, "sysRow")
	bankTrx := makeTrxMap("trx123", 3, "bankRow")

	count := deleteMatchedTrx(systemTrx, bankTrx)

	if count != 4 {
		t.Errorf("Expected 4 matched transactions, got %d", count)
	}
	if len(systemTrx) != 0 {
		t.Errorf("Expected system map to be empty after match")
	}
	if len(bankTrx["trx123"]) != 1 {
		t.Errorf("Expected 1 leftover in bank map, got %v", len(bankTrx["trx123"]))
	}
}

func TestDeleteMatchedTrx_MoreSystemTrx(t *testing.T) {
	systemTrx := makeTrxMap("trx123", 4, "sysRow")
	bankTrx := makeTrxMap("trx123", 3, "bankRow")

	count := deleteMatchedTrx(systemTrx, bankTrx)

	if count != 6 {
		t.Errorf("Expected 6 matched transactions, got %d", count)
	}
	if len(bankTrx) != 0 {
		t.Errorf("Expected bank map to be empty after match")
	}
	if len(systemTrx["trx123"]) != 1 {
		t.Errorf("Expected 1 leftover in system map, got %v", len(systemTrx["trx123"]))
	}
}
