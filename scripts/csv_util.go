package scripts

import (
	"encoding/csv"
	"os"
)

// WriteCSV writes the initial data to a CSV file without the GIF filenames.
func WriteCSV(filename string, data map[string][2]string) error {
	csvFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write CSV header
	err = csvWriter.Write([]string{"id", "signe", "definition"})
	if err != nil {
		return err
	}

	// Write CSV records
	for id, texts := range data {
		record := []string{id, texts[0], texts[1]}
		err := csvWriter.Write(record)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadCSV reads the CSV file and returns the records.
func ReadCSV(filename string) ([][]string, error) {
	csvFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer csvFile.Close()

	csvReader := csv.NewReader(csvFile)
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields per record
	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// WriteUpdatedCSV writes the updated records to a new CSV file.
func WriteUpdatedCSV(filename string, records [][]string) error {
	tmpCsvFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer tmpCsvFile.Close()

	tmpCsvWriter := csv.NewWriter(tmpCsvFile)
	defer tmpCsvWriter.Flush()

	err = tmpCsvWriter.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}
