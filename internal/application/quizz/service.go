package quizz

import (
	"encoding/csv"
	"log"
	"lsfb/internal/domain"
	"math/rand"
	"os"
)

var (
	typeSelection = "selection"
	typeSigne     = "signe"
	typeWrite     = "write"
)

var quizzes []domain.QuizInterface

type CSVRecord struct {
	ID         string
	Signe      string
	Definition string
}
type Service struct {
	records []CSVRecord
}

func NewQuizzService(csvFile string) *Service {
	s := &Service{}
	s.loadCSV(csvFile)
	return s
}

// TODO Replace with db
// loadCSV reads the CSV file and stores the data in the QuizService
func (s *Service) loadCSV(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rawRecords, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// Skipping the header row
	for _, rawRecord := range rawRecords[1:] {
		record := CSVRecord{
			ID:         rawRecord[0],
			Signe:      rawRecord[1],
			Definition: rawRecord[2],
		}
		s.records = append(s.records, record)
	}
}

// GetQuizz generates a quiz question
func (s *Service) GetQuizz() domain.QuizInterface {
	if len(s.records) < 4 {
		log.Fatal("Not enough records to generate a quiz question")
	}

	// Randomly select the correct answer
	correctIdx := rand.Intn(len(s.records))
	correctRecord := s.records[correctIdx]

	// Randomly select three other rows as distractors
	distractors := make(map[int]CSVRecord)
	for len(distractors) < 3 {
		idx := rand.Intn(len(s.records))
		if idx != correctIdx {
			distractors[idx] = s.records[idx]
		}
	}

	// Combine the correct answer and distractors
	signs := []string{correctRecord.Signe}
	for _, record := range distractors {
		signs = append(signs, record.Signe)
	}

	// Shuffle the signs and determine the new index of the correct answer
	shuffledSigns, _ := shuffleSigns(signs, correctRecord.Signe)

	return domain.SelectionQuiz{
		GifID: correctRecord.ID,
		Signe: shuffledSigns,
	}
}

// shuffleSigns shuffles the signs and returns the shuffled signs along with the new index of the correct sign
func shuffleSigns(signs []string, correctSign string) ([]string, int) {
	shuffledSigns := make([]string, len(signs))
	copy(shuffledSigns, signs)
	rand.Shuffle(len(shuffledSigns), func(i, j int) {
		shuffledSigns[i], shuffledSigns[j] = shuffledSigns[j], shuffledSigns[i]
	})

	var newCorrectIdx int
	for i, sign := range shuffledSigns {
		if sign == correctSign {
			newCorrectIdx = i
			break
		}
	}

	return shuffledSigns, newCorrectIdx
}
