package file

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

func WriteFileCSV(matrix [][]interface{}, folderName string, fileName string) {
	// Define relative path
	relativePath := filepath.Join(".", folderName, fileName)

	// Create directory if it doesn't exist
	err := os.MkdirAll(filepath.Dir(relativePath), os.ModePerm)
	if err != nil {
		fmt.Println("Error creating directory:", err)
		return
	}

	// Create or truncate the file
	file, err := os.Create(relativePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Loop through each row
	for _, row := range matrix {
		var stringRow []string
		for _, val := range row {
			stringRow = append(stringRow, fmt.Sprintf("%v", val)) // Convert each value to string
		}
		err := writer.Write(stringRow)
		if err != nil {
			fmt.Println("Error writing row to CSV:", err)
		}
	}

	fmt.Println("Transaction CSV file written to:", relativePath)
}
