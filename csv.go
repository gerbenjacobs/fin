package fin

import (
	"fmt"
	"os"
	"regexp"

	"github.com/gocarina/gocsv"
)

var currencyNormalizerRegex = regexp.MustCompile(`\([A-Z]{3}\)`)

func ReadCSVFile(fileName string, dst interface{}) error {
	csvFile, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("failed to load file %v", err)
	}
	defer csvFile.Close()

	gocsv.SetHeaderNormalizer(func(s string) string {
		return string(currencyNormalizerRegex.ReplaceAll([]byte(s), []byte("")))
	})
	return gocsv.UnmarshalFile(csvFile, dst)
}

func WriteCSVFile(fileName string, output interface{}) error {
	csvFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to load file %v", err)
	}

	return gocsv.MarshalFile(output, csvFile)
}
