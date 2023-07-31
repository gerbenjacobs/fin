package fin

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

func ReadCSVFile(fileName string, dst interface{}) error {
	csvFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to load file %v", err)
	}
	defer csvFile.Close()

	return gocsv.UnmarshalFile(csvFile, dst)
}

func WriteCSVFile(fileName string, output interface{}) error {
	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to load file %v", err)
	}

	return gocsv.MarshalFile(output, csvFile)
}
