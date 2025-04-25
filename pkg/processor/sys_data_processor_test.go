package processor

import (
	"reconciler/pkg/util"
	"testing"
	"time"
)

func TestProcessSystemTrxRecords(t *testing.T) {
	// Sample records for testing
	records := [][]string{
		{"T001", "100.00", "CREDIT", "2025-01-01 12:00:00"},
		{"T002", "50.50", "DEBIT", "2025-01-02 14:00:00"},
		{"T003", "200.75", "CREDIT", "2025-01-03 16:00:00"},
		{"T004", "invalid", "CREDIT", "2025-01-04 18:00:00"}, // Invalid amount
		{"T005", "150.00", "INVALID", "2025-01-05 20:00:00"}, // Invalid trx type
	}

	// Define filter dates
	filterStartDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	filterEndDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	// Create an empty result map to store processed transactions
	resultMap := make(map[string][]interface{})

	// Process the records
	count, totalAmount := ProcessSystemTrxRecords(records, resultMap, filterStartDate, filterEndDate)

	// Expected values after processing valid records
	expectedCount := 3            // 3 valid records
	expectedTotalAmount := 250.25 // Sum of valid amounts: 100.00 + (-50.50) + 200.75

	// Check that the correct number of transactions were processed
	if count != expectedCount {
		t.Errorf("expected %d transactions, got %d", expectedCount, count)
	}

	// Check the total amount (sum of valid transaction amounts)
	if totalAmount != expectedTotalAmount {
		t.Errorf("expected total amount %f, got %f", expectedTotalAmount, totalAmount)
	}

	// Check that the map has the correct keys and values for valid transactions
	if len(resultMap) != 3 { // 2 unique keys "CREDIT_100.0" and "DEBIT_50.5"
		t.Errorf("expected 2 unique keys, got %d", len(resultMap))
	}

	// Test if the key "CREDIT_100.00" exists and contains a transaction with amount 100.0
	cKey := "CREDIT_100.00"
	if _, exists := resultMap[cKey]; !exists {
		t.Errorf("expected key %s not found in resultMap", cKey)
	} else {
		// Check that the amount is correct
		if resultMap[cKey][0].(SystemTrx).AbsAmount != 100.0 {
			t.Errorf("expected amount 100.0 for key %s, got %f", cKey, resultMap[cKey][0].(SystemTrx).AbsAmount)
		}
	}

	// Test if the key "DEBIT_50.50" exists and contains a transaction with amount -50.5
	dKey := "DEBIT_50.50"
	if _, exists := resultMap[dKey]; !exists {
		t.Errorf("expected key %s not found in resultMap", dKey)
	} else {
		// Check that the amount is correct
		if resultMap[dKey][0].(SystemTrx).AbsAmount != 50.5 {
			t.Errorf("expected amount 50.5 for key %s, got %f", dKey, resultMap[dKey][0].(SystemTrx).AbsAmount)
		}
	}
}

func TestConstructSystemTrx(t *testing.T) {
	// Valid record for CREDIT
	row := []string{"T001", "100.00", "CREDIT", "2025-01-01 12:00:00"}
	filterStartDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	filterEndDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	mapKey, systemTrx, amount, err := ConstructSystemTrx(row, filterStartDate, filterEndDate)

	// Check if there's no error
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Validate the map key, system transaction, and amount
	if mapKey != "CREDIT_100.00" {
		t.Errorf("expected mapKey CREDIT_100.00, got %s", mapKey)
	}
	if systemTrx.TrxID != "T001" {
		t.Errorf("expected TrxID T001, got %s", systemTrx.TrxID)
	}
	if systemTrx.AbsAmount != 100.0 {
		t.Errorf("expected AbsAmount 100.0, got %f", systemTrx.AbsAmount)
	}
	if systemTrx.TrxType != util.TRX_TYPE_CREDIT {
		t.Errorf("expected TrxType CREDIT, got %s", systemTrx.TrxType)
	}
	expectedTrxTime := "2025-01-01 12:00:00"
	if systemTrx.TrxTime.Format("2006-01-02 15:04:05") != expectedTrxTime {
		t.Errorf("expected TrxTime %s, got %s", expectedTrxTime, systemTrx.TrxTime.Format("2006-01-02 15:04:05"))
	}
	if amount != 100.0 {
		t.Errorf("expected amount 100.0, got %f", amount)
	}

	// Invalid amount (should return an error)
	row = []string{"T002", "invalid", "CREDIT", "2025-01-02 14:00:00"}
	_, _, _, err = ConstructSystemTrx(row, filterStartDate, filterEndDate)
	if err == nil {
		t.Errorf("expected error for invalid amount, got nil")
	}

	// Invalid transaction type (should return an error)
	row = []string{"T003", "50.50", "INVALID", "2025-01-02 14:00:00"}
	_, _, _, err = ConstructSystemTrx(row, filterStartDate, filterEndDate)
	if err == nil {
		t.Errorf("expected error for invalid trx type, got nil")
	}
}
