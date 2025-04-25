package processor

import (
	"math"
	"path/filepath"
	"reconciler/pkg/logger"
	"reconciler/pkg/util"
	"regexp"
	"strconv"
	"time"
)

type BankTrx struct {
	TrxID    string
	Amount   float64
	TrxTime  time.Time
	BankName string
}

func ProcessBankTrxRecords(filePath string, records [][]string, mapResult map[string][]interface{}, filterStartDate time.Time, filterEndDate time.Time) (int, float64) {

	countTrx := 0
	totalAmount := float64(0)
	bankName := getBankName(filePath)

	for _, row := range records {

		mapKey, newBankTrx, amount, err := ConstructBankTrx(bankName, row, filterStartDate, filterEndDate)

		if err != nil {
			continue // skip this row
		}

		// put new bankTrx into map
		_, mapKeyExists := mapResult[mapKey]
		if mapKeyExists {
			mapResult[mapKey] = append(mapResult[mapKey], newBankTrx)
		} else {
			mapResult[mapKey] = []interface{}{newBankTrx}
		}

		// increase countTrx & totalAmount
		countTrx++
		totalAmount += amount
	}

	return countTrx, totalAmount
}

func ConstructBankTrx(bankName string, row []string, filterStartDate time.Time, filterEndDate time.Time) (string, BankTrx, float64, error) {
	trxID := row[0]

	colAmount := row[1]
	amount, err := strconv.ParseFloat(colAmount, 64)
	if err != nil {
		logger.Warn("%s : %v", "Invalid amount: ", colAmount)
		// skip this row
		return "", BankTrx{}, 0, err
	}

	colTrxDate := row[2]
	trxDate, err := time.Parse("2006-01-02", colTrxDate)
	if err != nil {
		logger.Warn("%s : %v", "Invalid trxDate: ", colTrxDate)
		// skip this row
		return "", BankTrx{}, 0, err
	}

	if !util.IsWithinRange(trxDate, filterStartDate, filterEndDate) {
		logger.Info("%s : %v", "Skip this trx on: ", trxDate)

		// skip this row
		return "", BankTrx{}, 0, err
	}

	trxType := util.TRX_TYPE_CREDIT
	if amount < 0 {
		trxType = util.TRX_TYPE_DEBIT
	}
	mapKey := ConstructMapKey(trxType, math.Abs(amount))
	newBankTrx := BankTrx{BankName: bankName, TrxID: trxID, Amount: amount, TrxTime: trxDate}

	return mapKey, newBankTrx, amount, nil
}

func getBankName(path string) string {
	fileName := filepath.Base(path)
	re := regexp.MustCompile(`bank_[a-zA-Z0-9]*`)
	bankName := re.FindString(fileName)

	return bankName
}

func GroupBankTrx(bankTrxMap map[string][]interface{}) map[string][]interface{} {

	if len(bankTrxMap) == 0 {
		return nil
	}

	// grouped by bankName
	groupedBankTrx := make(map[string][]interface{})

	for _, bankTrxList := range bankTrxMap {

		for _, bankTrxInterface := range bankTrxList {

			bankTrx := bankTrxInterface.(BankTrx)
			bankName := bankTrx.BankName

			// copy the bankTrx to groupedBankTrx
			_, mapKeyExists := groupedBankTrx[bankName]
			if mapKeyExists {
				groupedBankTrx[bankName] = append(groupedBankTrx[bankName], bankTrx)
			} else {
				groupedBankTrx[bankName] = []interface{}{bankTrx}
			}
		}

	}

	return groupedBankTrx
}
