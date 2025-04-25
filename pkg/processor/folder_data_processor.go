package processor

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reconciler/pkg/file"
	"reconciler/pkg/util"
	"time"
)

// ProcessFolderCSV read all .csv files in a folder and returns map data
func ProcessFolderCSV(folderPath string, folderType string, filterStartDate time.Time, filterEndDate time.Time) (map[string][]interface{}, int, float64) {

	resultMap := make(map[string][]interface{})
	countTrx := 0
	totalAmount := float64(0)

	// Walk through the folder and find CSV files
	err := filepath.WalkDir(folderPath, func(path string, d fs.DirEntry, err error) error {

		if d.IsDir() {
			return nil // This is only directory description. Skip this WalkDir iteration.
		}

		records := file.ReadFileCSV(path)
		if records == nil {
			return nil
		}

		countTrxPerFile := 0
		totalAmountPerFile := float64(0)

		if folderType == util.FILE_TYPE_BANK_TRX {
			countTrxPerFile, totalAmountPerFile = ProcessBankTrxRecords(path, records, resultMap, filterStartDate, filterEndDate)

		} else if folderType == util.FILE_TYPE_SYSTEM_TRX {
			countTrxPerFile, totalAmountPerFile = ProcessSystemTrxRecords(records, resultMap, filterStartDate, filterEndDate)
		}

		// increase transaction count & total amount
		countTrx += countTrxPerFile
		totalAmount += totalAmountPerFile

		return nil
	})

	if err != nil {
		return nil, 0, float64(0)
	}

	return resultMap, countTrx, totalAmount
}

func ConstructMapKey(trxType string, absAmount float64) string {
	strAbsAmount := fmt.Sprintf("%.2f", absAmount)
	mapKey := trxType + "_" + strAbsAmount
	return mapKey
}
