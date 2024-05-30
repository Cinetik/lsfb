package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/schollz/progressbar/v3"
)

func main() {
	// Create a directory to store GIFs
	err := os.Mkdir("gif", 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Scrape data from the website
	idMap, err := ScrapeData()
	if err != nil {
		log.Fatalf("Failed to scrape data: %v", err)
	}

	// Write initial CSV without GIF filenames
	err = WriteCSV("output.csv", idMap)
	if err != nil {
		log.Fatalf("Failed to write CSV file: %v", err)
	}

	// Read the CSV content
	records, err := ReadCSV("output.csv")
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
	}

	// Initialize progress bar
	totalIDs := len(records) - 1 // Subtracting header row
	bar := progressbar.Default(int64(totalIDs))

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Loop through each record and download GIFs
	for i, record := range records {
		// Skip header row
		if i == 0 {
			continue
		}

		wg.Add(1) // Increment the WaitGroup counter for each goroutine

		id := record[0]
		signe := record[1]
		definition := record[2]

		// Goroutine to download GIFs
		go func(id, signe, definition string, rowIndex int) {
			defer wg.Done() // Decrement the WaitGroup counter when the goroutine completes

			// Download GIF and update records
			err := DownloadGIFAndUpdateRecords(id, signe, definition, rowIndex, records)
			if err != nil {
				log.Printf("Error downloading GIF for ID %s: %v", id, err)
			}

			// Update progress bar
			bar.Add(1)
		}(id, signe, definition, i) // Pass the ID, signe, definition, and row index to the goroutine

		// Introduce a delay between goroutines to avoid overloading the server
		time.Sleep(100 * time.Millisecond)
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Write updated records to a new CSV file
	err = WriteUpdatedCSV("output_tmp.csv", records)
	if err != nil {
		log.Fatalf("Failed to write updated records to temporary CSV file: %v", err)
	}

	// Replace the original CSV file with the updated one
	err = os.Rename("output_tmp.csv", "output.csv")
	if err != nil {
		log.Fatalf("Failed to replace the original CSV file: %v", err)
	}

	fmt.Println("CSV file and GIFs downloaded successfully.")
}
