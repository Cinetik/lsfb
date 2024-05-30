package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/schollz/progressbar/v3"
)

// GIFResponse represents the JSON response from the GIF query
type GIFResponse struct {
	URL string `json:"url"`
}

func main() {
	// Base URL for the lexique page
	baseURL := "https://www.corpus-lsfb.be/lexique.php?lettre="

	// Map to store the extracted data
	idMap := make(map[string][2]string)

	// Loop through letters a to z
	for letter := 'a'; letter <= 'z'; letter++ {
		// Construct the URL with the current letter
		url := baseURL + string(letter)

		// Make HTTP request with retries and exponential backoff
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to retrieve data for letter %c: %v", letter, err)
			continue
		}
		defer resp.Body.Close()

		// Check if the request was successful
		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to retrieve data for letter %c: received status code %d", letter, resp.StatusCode)
			continue
		}

		// Parse the HTML content with goquery
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Printf("Failed to parse HTML for letter %c: %v", letter, err)
			continue
		}

		// Find all tr tags with the 'vocabulaire' class
		doc.Find("tr.vocabulaire").Each(func(index int, item *goquery.Selection) {
			// Extract the id attribute
			if id, exists := item.Attr("id"); exists {
				// Get the first and third td tags' text
				firstTD := item.Find("td").Eq(0).Text()
				thirdTD := item.Find("td").Eq(2).Text()
				idMap[id] = [2]string{firstTD, thirdTD}
			}
		})

		// Introduce a delay between requests to avoid overloading the server
		time.Sleep(500 * time.Millisecond)
	}

	// Create a CSV file without the GIF filenames
	csvFile, err := os.Create("output.csv")
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer csvFile.Close()

	// Write to CSV file
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	// Write CSV header
	err = csvWriter.Write([]string{"id", "signe", "definition", "gif"})
	if err != nil {
		log.Fatalf("Failed to write CSV header: %v", err)
	}

	// Create a directory to store GIFs
	err = os.Mkdir("gif", 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create directory: %v", err)
	}

	// Initialize progress bar
	totalIDs := len(idMap)
	bar := progressbar.Default(int64(totalIDs))

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Loop through each ID in the map
	for id, texts := range idMap {
		wg.Add(1) // Increment the WaitGroup counter for each goroutine

		// Write CSV record
		record := []string{id, texts[0], texts[1]}
		err := csvWriter.Write(record)
		if err != nil {
			log.Fatalf("Failed to write CSV record: %v", err)
		}

		// Goroutine to download GIFs
		go func(id, signe, definition string) {
			defer wg.Done() // Decrement the WaitGroup counter when the goroutine completes

			// Download GIF and update CSV
			err := downloadGIFAndUpdateCSV(id)
			if err != nil {
				log.Printf("Error downloading GIF for ID %s: %v", id, err)
			}
		}(id, texts[0], texts[1]) // Pass the ID, signe, and definition to the goroutine

		// Introduce a delay between goroutines to avoid overloading the server
		time.Sleep(100 * time.Millisecond)
	}

	// Goroutine to wait for all other goroutines to finish
	go func() {
		wg.Wait() // Wait for all goroutines to finish
	}()

	// Update progress bar
	bar.Finish()

	fmt.Println("CSV file and GIFs downloaded successfully.")
}

// downloadGIFAndUpdateCSV downloads the GIF for the given ID and updates the CSV with the filename
func downloadGIFAndUpdateCSV(id string) error {
	// Query URL for the GIF
	gifURL := "https://www.corpus-lsfb.be/getVocabulaire.php?mot=" + id

	// Make a GET request to fetch the GIF JSON response
	resp, err := http.Get(gifURL)
	if err != nil {
		return fmt.Errorf("failed to retrieve GIF for ID %s: %v", id, err)
	}
	defer resp.Body.Close()

	// Check if the request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to retrieve GIF for ID %s: received status code %d", id, resp.StatusCode)
	}

	// Parse the JSON response
	var gifResponse GIFResponse
	err = json.NewDecoder(resp.Body).Decode(&gifResponse)
	if err != nil {
		return fmt.Errorf("failed to parse GIF response for ID %s: %v", id, err)
	}

	// Extract filename from URL
	filename := filepath.Base(gifResponse.URL)

	// Create a file to save the GIF
	gifPath := filepath.Join("gif", filename)

	// Check if the file already exists locally and has a size greater than 0KB
	if fileInfo, err := os.Stat(gifPath); err == nil && fileInfo.Size() > 0 {
		log.Printf("Skipping download for ID %s, GIF already exists locally", id)
		return nil
	}

	// Create a file to save the GIF
	file, err := os.Create(gifPath)
	if err != nil {
		return fmt.Errorf("failed to create GIF file for ID %s: %v", id, err)
	}
	defer file.Close()

	// Make a GET request to download the GIF
	resp, err = http.Get("https://www.corpus-lsfb.be/img/pictures/" + gifResponse.URL)
	if err != nil {
		return fmt.Errorf("failed to download GIF for ID %s: %v", id, err)
	}
	defer resp.Body.Close()

	// Copy the GIF content to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save GIF for ID %s: %v", id, err)
	}

	// Update CSV with the GIF filename
	err = updateCSVWithGIF(id, filename)
	if err != nil {
		return fmt.Errorf("failed to update CSV with GIF filename for ID %s: %v", id, err)
	}

	return nil
}

// updateCSVWithGIF updates the CSV with the GIF filename for the given ID
func updateCSVWithGIF(id, filename string) error {
	// Open the CSV file
	csvFile, err := os.OpenFile("output.csv", os.O_RDWR, 0644)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %v", err)
	}
	defer csvFile.Close()

	// Read the CSV content
	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV content: %v", err)
	}

	// Find the row corresponding to the ID and update the GIF filename
	for i, record := range records {
		if record[0] == id {
			records[i] = append(record, filename)
			break
		}
	}

	// Write the updated CSV content
	csvFile.Seek(0, 0)
	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()
	err = csvWriter.WriteAll(records)
	if err != nil {
		return fmt.Errorf("failed to write updated CSV content: %v", err)
	}

	return nil
}
