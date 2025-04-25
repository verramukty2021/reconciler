package processor

import (
	"testing"
	"time"
)

func TestGetBankName(t *testing.T) {
	filePath := "/some/path/bank_bca_data.csv"
	bank := getBankName(filePath)
	if bank != "bank_bca" {
		t.Errorf("Expected bank name 'bank_bca', got '%s'", bank)
	}
}

func TestConstructBankTrx(t *testing.T) {
	row := []string{"TRX001", "-250.75", "2025-01-15"}
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	mapKey, trx, amt, err := ConstructBankTrx("bank_bca", row, start, end)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if trx.BankName != "bank_bca" || trx.TrxID != "TRX001" {
		t.Errorf("Unexpected transaction data: %+v", trx)
	}
	if amt != -250.75 {
		t.Errorf("Expected amount -250.75, got %f", amt)
	}
	if mapKey != "DEBIT_250.75" {
		t.Errorf("Unexpected map key: %s", mapKey)
	}
}

func TestProcessBankTrxRecords(t *testing.T) {
	records := [][]string{
		{"TRX001", "100.00", "2025-01-15"},
		{"TRX002", "-50.00", "2025-01-20"},
	}
	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	mapResult := make(map[string][]interface{})
	count, total := ProcessBankTrxRecords("bank_bca_file.csv", records, mapResult, start, end)

	if count != 2 {
		t.Errorf("Expected 2 transactions, got %d", count)
	}
	if total != 50.00 {
		t.Errorf("Expected total amount 50.00, got %.2f", total)
	}
	if len(mapResult) != 2 {
		t.Errorf("Expected 2 keys in map, got %d", len(mapResult))
	}
}

func TestGroupBankTrx(t *testing.T) {
	trx1 := BankTrx{TrxID: "1", Amount: 100.0, TrxTime: time.Now(), BankName: "bank_bca"}
	trx2 := BankTrx{TrxID: "2", Amount: 200.0, TrxTime: time.Now(), BankName: "bank_mandiri"}

	inputMap := map[string][]interface{}{
		"C-100.00": {trx1},
		"C-200.00": {trx2},
	}

	grouped := GroupBankTrx(inputMap)

	if len(grouped) != 2 {
		t.Errorf("Expected 2 grouped bank names, got %d", len(grouped))
	}
	if len(grouped["bank_bca"]) != 1 {
		t.Errorf("Expected 1 trx for bank_bca")
	}
	if len(grouped["bank_mandiri"]) != 1 {
		t.Errorf("Expected 1 trx for bank_mandiri")
	}
}
