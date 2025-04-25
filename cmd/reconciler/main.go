package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"reconciler/pkg/processor"
	"reconciler/pkg/util"
	"time"
)

func main() {

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	// Get input from command
	// os.Args[0] is the program name
	args := os.Args[1:] // skip the program name

	if len(args) > 0 {
		startDate, endDate = readInputFromCommand(args)
	}

	// process input csv
	bankFolderPath := filepath.Join(".", "csv/input_bank_trx")
	sysTrxFolderPath := filepath.Join(".", "csv/input_system_trx")
	bankTrxMap, countAllBankTrx, bankTotalAmount := processor.ProcessFolderCSV(bankFolderPath, util.FILE_TYPE_BANK_TRX, startDate, endDate)
	systemTrxMap, countAllSysTrx, sysTotalAmount := processor.ProcessFolderCSV(sysTrxFolderPath, util.FILE_TYPE_SYSTEM_TRX, startDate, endDate)

	// delete matched trx on both map (sysTrxMap & bankTrxMap)
	countAllMatchedTrx := deleteMatchedTrx(systemTrxMap, bankTrxMap)

	// Group the leftover unmatch bankTrx by Bank.
	unmatchedBankTrx := processor.GroupBankTrx(bankTrxMap)

	// count all metric
	countAllTrx := countAllSysTrx + countAllBankTrx
	countAllUnmatch := countAllTrx - countAllMatchedTrx

	// Notes: totalAmountDiscrepancy can also be calculated using unmatched total amount
	// The alternative formula: math.Abs(sysTotalAmountUnmatched - bankTotalAmountUnmatched)
	// but we use the overall total amount because we already have the value during initial csv file processing.
	totalAmountDiscrepancy := util.RoundToTwoDecimals(math.Abs(sysTotalAmount - bankTotalAmount))

	// print result
	printResult(countAllTrx, countAllMatchedTrx, countAllUnmatch, systemTrxMap, unmatchedBankTrx, totalAmountDiscrepancy)

}

func readInputFromCommand(args []string) (time.Time, time.Time) {
	cmdMessage := "Usage: go run cmd/reconciler/main.go <startDate> <endDate> \n"
	cmdMessage += "Example: go run cmd/reconciler/main.go 2025-01-01 2025-01-31 \n"

	if len(args) != 2 {
		fmt.Printf("%s\n", cmdMessage)
		return time.Time{}, time.Time{}
	}

	startDate, err := time.Parse("2006-01-02", args[0])
	endDate, err := time.Parse("2006-01-02", args[1])

	if err != nil {
		fmt.Printf("%s\n %v \n", cmdMessage, err.Error())
		return time.Time{}, time.Time{}
	}
	return startDate, endDate
}

func printResult(countAllTrx int, countAllMatch int, countAllUnmatch int, systemTrxMap map[string][]interface{}, unmatchedBankTrx map[string][]interface{}, totalAmountDiscrepancy float64) {
	fmt.Printf("\n%s: \n%v\n", "Total number of transactions processed", countAllTrx)
	fmt.Printf("\n%s: \n%v\n", "Total number of matched transactions", countAllMatch)
	fmt.Printf("\n%s: \n%v\n", "Total number of unmatched transactions", countAllUnmatch)

	if countAllUnmatch > 0 {
		fmt.Printf("\n%s \n%s: \n%s \n", "---------------", "Details of unmatched transactions", "---------------")
	}

	if len(systemTrxMap) > 0 {
		fmt.Printf("\n%s: \n", "System transaction details missing in bank statement(s)")

		// There's only unmatch sys trx left on systemTrxMap
		for _, unmatchSysTrx := range systemTrxMap {
			for _, row := range unmatchSysTrx {
				fmt.Printf("%v\n", row)
			}
		}
	}

	if unmatchedBankTrx != nil || len(unmatchedBankTrx) > 0 {
		fmt.Printf("\n%s: \n", "Bank statement details missing in system transactions (grouped by bank)")

		for bankName, bankTrx := range unmatchedBankTrx {
			fmt.Printf("\n%s: %v\n", "Bank Name", bankName)

			for _, row := range bankTrx {
				fmt.Printf("%v\n", row)
			}
		}

	}

	fmt.Printf("\n%s: \n%v\n", "Total discrepancies (sum of absolute differences in amount between matched transactions)", totalAmountDiscrepancy)
}

func deleteMatchedTrx(systemTrxMap map[string][]interface{}, bankTrxMap map[string][]interface{}) int {

	countAllMatchedTrx := 0

	for key, systemTrx := range systemTrxMap {

		// check whether the system trx exists on bank trx
		bankTrx, bankTrxExists := bankTrxMap[key]

		if bankTrxExists {

			countSysTrx := len(systemTrx)
			countBankTrx := len(bankTrx)

			if countSysTrx == countBankTrx {
				// all system trx match on bank trx data
				// delete all from both map(sysTrx & bankTrx) for this key
				delete(systemTrxMap, key)
				delete(bankTrxMap, key)

				// increase match
				countAllMatchedTrx += countSysTrx + countBankTrx

			} else if countSysTrx < countBankTrx {
				// there are some bank trx that's missing on system trx
				// delete all from sysTrx for this key
				delete(systemTrxMap, key)

				// delete partial from bankTrx
				isDeleted, err := util.DeleteFromSliceMap(bankTrxMap, key, 0, countSysTrx)
				if err != nil || !isDeleted {
					fmt.Printf("%s %v \n", "Failure on DeleteFromSliceMap", err)
					os.Exit(1) // stop application, because the failure on this process can cause wrong calculation.
				}

				// increase match with countSysTrx
				countAllMatchedTrx += countSysTrx * 2

			} else {
				// totalSystemTrx > totalBankTrx
				// there are systemTrx that's missing on bankTrx
				// delete all from bankTrx for this key
				delete(bankTrxMap, key)

				// delete partial from sysTrx
				isDeleted, err := util.DeleteFromSliceMap(systemTrxMap, key, 0, countBankTrx)
				if err != nil || !isDeleted {
					fmt.Printf("%s %v \n", "Failure on DeleteFromSliceMap", err)
					os.Exit(1) // stop application, because the failure on this process can cause wrong calculation.
				}

				// increase match with countBankTrx
				countAllMatchedTrx += countBankTrx * 2
			}

		}
	}
	return countAllMatchedTrx
}
