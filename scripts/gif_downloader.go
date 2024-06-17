package scripts

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// GIFResponse represents the JSON response from the GIF query
type GIFResponse struct {
	URL string `json:"url"`
}

var mu sync.Mutex

// DownloadGIFAndUpdateRecords downloads the GIF for the given ID and updates the records with the filename.
func DownloadGIFAndUpdateRecords(id, signe, definition string, rowIndex int, records [][]string) error {
	// Construct the filename for the GIF
	filename := id + ".gif"
	filepath := filepath.Join("gif", filename)

	// Check if the file already exists locally and has a size greater than 0KB
	if fileInfo, err := os.Stat(filepath); err == nil && fileInfo.Size() > 0 {
		log.Printf("Skipping download for ID %s, GIF already exists locally", id)
		mu.Lock()
		records[rowIndex] = append(records[rowIndex], filename)
		mu.Unlock()
		return nil
	}

	// Retry mechanism with exponential backoff
	var resp *http.Response
	var err error
	for i := 0; i < 5; i++ {
		resp, err = http.Get("https://www.corpus-lsfb.be/getVocabulaire.php?mot=" + id)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(1<<i) * time.Second) // Exponential backoff
	}
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

	// Create a file to save the GIF
	file, err := os.Create(filepath)
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

	// Update records with the GIF filename
	mu.Lock()
	records[rowIndex] = append(records[rowIndex], filename)
	mu.Unlock()

	return nil
}
