package processor

import (
	"os"
	"path/filepath"
	"reconciler/pkg/util"
	"testing"
	"time"
)

// Test ProcessFolderCSV function
func TestProcessFolderCSV(t *testing.T) {

	// Create a temporary folder for the test
	tmpFolder := "testdata"
	err := os.MkdirAll(tmpFolder, os.ModePerm)
	if err != nil {
		t.Fatalf("failed to create test folder: %v", err)
	}
	defer os.RemoveAll(tmpFolder)

	// Add some CSV files to the folder for testing
	bankFile := filepath.Join(tmpFolder, "bank_trx.csv")
	systemFile := filepath.Join(tmpFolder, "system_trx.csv")

	// Create mock CSV files for testing purposes
	_ = os.WriteFile(bankFile, []byte("bank_trx_id,amount,bank_trx_date\nT001,100.00,2025-01-01\nT002,-50.50,2025-01-02"), os.ModePerm)
	_ = os.WriteFile(systemFile, []byte("trx_id,amount,type,trx_datetime\nT003,200.75,CREDIT,2025-01-03 16:00:00\nT004,150.00,DEBIT,2025-01-04 18:00:00"), os.ModePerm)

	// Define filter dates
	filterStartDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	filterEndDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	// Call the ProcessFolderCSV function
	resultMap, countTrx, totalAmount := ProcessFolderCSV(tmpFolder, util.FILE_TYPE_BANK_TRX, filterStartDate, filterEndDate)

	// Test the expected results after processing
	expectedCount := 2
	expectedTotalAmount := 49.5 // Sum of the valid transaction amounts
	if countTrx != expectedCount {
		t.Errorf("expected %d transactions, got %d", expectedCount, countTrx)
	}

	if totalAmount != expectedTotalAmount {
		t.Errorf("expected total amount %f, got %f", expectedTotalAmount, totalAmount)
	}

	// Check the result map for the expected keys and values
	expectedKey := "CREDIT_100.00"
	if _, exists := resultMap[expectedKey]; !exists {
		t.Errorf("expected key %s not found in resultMap", expectedKey)
	}

	// Test for SYSTEM_TRX folder type
	resultMap, countTrx, totalAmount = ProcessFolderCSV(tmpFolder, util.FILE_TYPE_SYSTEM_TRX, filterStartDate, filterEndDate)

	expectedCount = 2
	expectedTotalAmount = 50.75 // Sum of valid amounts: 200.75 + (-150.00)
	if countTrx != expectedCount {
		t.Errorf("expected %d transactions, got %d", expectedCount, countTrx)
	}

	if totalAmount != expectedTotalAmount {
		t.Errorf("expected total amount %f, got %f", expectedTotalAmount, totalAmount)
	}

}

// Test ConstructMapKey function
func TestConstructMapKey(t *testing.T) {
	tests := []struct {
		key       string
		absAmount float64
		expected  string
	}{
		{"CREDIT", 100.00, "CREDIT_100.00"},
		{"CREDIT", 100.0, "CREDIT_100.00"},
		{"DEBIT", 100, "DEBIT_100.00"},
	}

	for _, tt := range tests {
		result := ConstructMapKey(tt.key, tt.absAmount)
		if result != tt.expected {
			t.Errorf("ConstructMapKey(%v, %v) = %v; want %v", tt.key, tt.absAmount, result, tt.expected)
		}
	}

}
