package file

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func ReadFileCSV(path string) [][]string {

	fileErrorFmt := "%s: %v\n %v\n \n"

	if filepath.Ext(path) != ".csv" {
		fmt.Printf(fileErrorFmt, "File is not CSV", path, nil)
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf(fileErrorFmt, "Cannot open file", path, err)
		return nil
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf(fileErrorFmt, "Cannot read file", path, err)
		return nil
	}

	// remove header row
	records = records[1:]

	if len(records) < 1 {
		fmt.Printf(fileErrorFmt, "File is empty", path, err)
		return nil
	}

	return records
}
