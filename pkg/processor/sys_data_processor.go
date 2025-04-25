package processor

import (
	"errors"
	"math"
	"reconciler/pkg/logger"
	"reconciler/pkg/util"
	"strconv"
	"time"
)

type SystemTrx struct {
	TrxID     string
	AbsAmount float64
	TrxType   string
	TrxTime   time.Time
}

func ProcessSystemTrxRecords(records [][]string, resultMap map[string][]interface{}, filterStartDate time.Time, filterEndDate time.Time) (int, float64) {

	countTrx := 0
	totalAmount := float64(0)

	for _, row := range records {

		mapKey, newSystemTrx, amount, err := ConstructSystemTrx(row, filterStartDate, filterEndDate)

		if err != nil {
			continue // skip this row
		}

		_, mapKeyExists := resultMap[mapKey]
		if mapKeyExists {
			resultMap[mapKey] = append(resultMap[mapKey], newSystemTrx)
		} else {
			resultMap[mapKey] = []interface{}{newSystemTrx}
		}

		// increase countTrx & totalAmount
		countTrx++
		totalAmount += amount
	}

	return countTrx, totalAmount
}

func ConstructSystemTrx(row []string, filterStartDate time.Time, filterEndDate time.Time) (string, SystemTrx, float64, error) {
	trxID := row[0]

	trxType := row[2]
	if trxType != util.TRX_TYPE_DEBIT && trxType != util.TRX_TYPE_CREDIT {
		return "", SystemTrx{}, 0, errors.New("invalid trx type")
	}

	absAmount, err := strconv.ParseFloat(row[1], 64)
	if err != nil {
		return "", SystemTrx{}, 0, err
	}

	colTrxDateTime := row[3]
	trxDateTime, err := time.Parse("2006-01-02 15:04:05", colTrxDateTime)
	if err != nil {
		logger.Warn("%s : %v", "Invalid trxDateTime: ", colTrxDateTime)

		// skip this row
		return "", SystemTrx{}, 0, err
	}

	if !util.IsWithinRange(trxDateTime, filterStartDate, filterEndDate) {
		logger.Info("%s : %v", "Skip this trx on: ", trxDateTime)

		// skip this row
		return "", SystemTrx{}, 0, err
	}

	mapKey := ConstructMapKey(trxType, absAmount)
	newSystemTrx := SystemTrx{TrxID: trxID, AbsAmount: absAmount, TrxType: trxType, TrxTime: trxDateTime}

	// get the real amount based on trxType
	realAmount := absAmount
	if trxType == util.TRX_TYPE_DEBIT {
		realAmount = -math.Abs(absAmount)
	}

	return mapKey, newSystemTrx, realAmount, nil
}
