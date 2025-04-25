package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"reconciler/pkg/dummy"
	"reconciler/pkg/file"
	"reconciler/pkg/util"
	"strconv"
	"time"
)

func main() {

	// Default input
	totalTrx := 10
	totalBankA := 3
	totalBankB := 3
	totalBankC := 4
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	// Get input from command
	// os.Args[0] is the program name
	args := os.Args[1:] // skip the program name

	if len(args) > 0 {
		totalTrx, totalBankA, totalBankB, totalBankC, startDate, endDate = readInputFromCommand(args)
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// init data with header
	trxData := [][]interface{}{
		{"trx_id", "amount", "type", "trx_datetime"},
	}
	bankDataA := [][]interface{}{
		{"bank_trx_id", "amount", "bank_trx_date"},
	}
	bankDataB := [][]interface{}{
		{"bank_trx_id", "amount", "bank_trx_date"},
	}
	bankDataC := [][]interface{}{
		{"bank_trx_id", "amount", "bank_trx_date"},
	}

	endOfBankA := totalBankA
	endOfBankB := endOfBankA + totalBankB
	endOfBankC := endOfBankB + totalBankC

	for i := 0; i < totalTrx; i++ {
		newTrxData, id, amount, trxType := generateNewTrxData(startDate, endDate)
		trxData = append(trxData, newTrxData)

		// split trxData into bank data
		newBankTrxData := generateNewBankData(startDate, endDate, id, amount, trxType)
		if i < endOfBankA {
			bankDataA = append(bankDataA, newBankTrxData)
		} else if i < endOfBankB {
			bankDataB = append(bankDataB, newBankTrxData)
		} else if i < endOfBankC {
			bankDataC = append(bankDataC, newBankTrxData)
		}

	}

	// write csv files
	file.WriteFileCSV(trxData, "csv/input_system_trx", "system_trx.csv")
	file.WriteFileCSV(bankDataA, "csv/input_bank_trx", "bank_a_trx.csv")
	file.WriteFileCSV(bankDataB, "csv/input_bank_trx", "bank_b_trx.csv")
	file.WriteFileCSV(bankDataC, "csv/input_bank_trx", "bank_c_trx.csv")
}

func readInputFromCommand(args []string) (int, int, int, int, time.Time, time.Time) {
	cmdMessage := "Usage: go run cmd/csv_generator/main.go <totalTrx> <totalBankA> <totalBankB> <totalBankB> <startDate> <endDate> \n"
	cmdMessage += "Example: go run cmd/csv_generator/main.go 10 3 3 4 2025-01-01 2025-01-31 \n"

	if len(args) != 6 {
		fmt.Printf("%s\n", cmdMessage)
		return 0, 0, 0, 0, time.Time{}, time.Time{}
	}

	totalTrx, err := strconv.Atoi(args[0])
	totalBankA, err := strconv.Atoi(args[1])
	totalBankB, err := strconv.Atoi(args[2])
	totalBankC, err := strconv.Atoi(args[3])
	startDate, err := time.Parse("2006-01-02", args[4])
	endDate, err := time.Parse("2006-01-02", args[5])

	if err != nil {
		fmt.Printf("%s\n %v \n", cmdMessage, err.Error())
		return 0, 0, 0, 0, time.Time{}, time.Time{}
	}
	return totalTrx, totalBankA, totalBankB, totalBankC, startDate, endDate
}

func generateNewBankData(startDate time.Time, endDate time.Time, id string, amount float64, trxType string) []interface{} {
	bankTrxDate := dummy.GenerateRandomTime(startDate, endDate).Format("2006-01-02") // Format as YYYY-MM-

	bankAmount := amount
	if trxType == util.TRX_TYPE_DEBIT {
		bankAmount = -math.Abs(amount)
	}

	newBankTrxData := []interface{}{id, bankAmount, bankTrxDate}
	return newBankTrxData
}

func generateNewTrxData(startDate time.Time, endDate time.Time) ([]interface{}, string, float64, string) {
	id := dummy.GenerateRandomID()

	const minAmount = 1000.0
	const maxAmount = 5000.0
	amount := dummy.GenerateRandomAmount(minAmount, maxAmount)

	trxType := dummy.GenerateRandomTrxType()
	trxDateTime := dummy.GenerateRandomTime(startDate, endDate).Format("2006-01-02 15:04:05") // Format: YYYY-MM-DD HH:MM:SS

	newTrxData := []interface{}{id, amount, trxType, trxDateTime}

	return newTrxData, id, amount, trxType
}
