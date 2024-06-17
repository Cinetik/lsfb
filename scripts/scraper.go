package scripts

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// ScrapeData scrapes data from the website and returns a map of IDs to their corresponding signe and definition.
func ScrapeData() (map[string][2]string, error) {
	// Base URL for the lexique page
	baseURL := "https://www.corpus-lsfb.be/lexique.php?lettre="

	// Map to store the extracted data
	idMap := make(map[string][2]string)

	// Loop through letters a to z
	for letter := 'a'; letter <= 'z'; letter++ {
		// Construct the URL with the current letter
		url := baseURL + string(letter)

		// Make HTTP request
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
				firstTD := strings.TrimSpace(item.Find("td").Eq(0).Text())
				thirdTD := strings.TrimSpace(item.Find("td").Eq(2).Text())
				idMap[id] = [2]string{firstTD, thirdTD}
			}
		})

		// Introduce a delay between requests to avoid overloading the server
		time.Sleep(500 * time.Millisecond)
	}

	return idMap, nil
}
