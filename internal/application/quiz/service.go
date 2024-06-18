package quiz

import (
	"encoding/csv"
	"log"
	"lsfb/internal/domain"
	"math/rand"
	"os"
)

var (
	typeSelection = "selection"
	typeSign      = "signe"
	typeWrite     = "write"
)

var quizzes []domain.QuizQuestionType

type CSVRecord struct {
	ID         string
	Sign       string
	Definition string
}
type Service struct {
	records []CSVRecord
}

func NewQuizService(csvFile string) *Service {
	s := &Service{}
	s.loadCSV(csvFile)
	return s
}

// TODO Replace with db
// loadCSV reads the CSV file and stores the data in the Service
func (qs *Service) loadCSV(filename string) {
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
			Sign:       rawRecord[1],
			Definition: rawRecord[2],
		}
		qs.records = append(qs.records, record)
	}
}
func (qs *Service) GetQuiz() domain.QuizQuestionType {
	switch rand.Intn(3) {
	case 0:
		return qs.generateQuizQuestion()
	case 1:
		return qs.generateDefinitionQuizQuestion()
	default:
		return qs.generateWriteDefinitionQuizQuestion()
	}
}

// generateQuizQuestion generates a quiz question with signs
func (qs *Service) generateQuizQuestion() *domain.QuizQuestion {
	if len(qs.records) < 4 {
		log.Fatal("Not enough records to generate a quiz question")
	}

	// Randomly select the correct answer
	correctIdx := rand.Intn(len(qs.records))
	correctRecord := qs.records[correctIdx]

	// Randomly select three other rows as distractors
	distractors := make(map[int]CSVRecord)
	for len(distractors) < 3 {
		idx := rand.Intn(len(qs.records))
		if idx != correctIdx {
			distractors[idx] = qs.records[idx]
		}
	}

	// Combine the correct answer and distractors
	signs := []string{correctRecord.Sign}
	for _, record := range distractors {
		signs = append(signs, record.Sign)
	}

	// Shuffle the signs and determine the new index of the correct answer
	shuffledSigns, newCorrectIdx := shuffleSigns(signs, correctRecord.Sign)

	return &domain.QuizQuestion{
		ID:         correctRecord.ID,
		Signs:      shuffledSigns,
		CorrectIdx: newCorrectIdx,
	}
}

// generateDefinitionQuizQuestion generates a quiz question with definitions
func (qs *Service) generateDefinitionQuizQuestion() *domain.DefinitionQuizQuestion {
	if len(qs.records) < 4 {
		log.Fatal("Not enough records to generate a quiz question")
	}

	// Randomly select the correct answer
	correctIdx := rand.Intn(len(qs.records))
	correctRecord := qs.records[correctIdx]

	// Randomly select three other rows as distractors
	distractors := make(map[int]CSVRecord)
	for len(distractors) < 3 {
		idx := rand.Intn(len(qs.records))
		if idx != correctIdx {
			distractors[idx] = qs.records[idx]
		}
	}

	// Combine the correct answer and distractors
	ids := []string{correctRecord.ID}
	for _, record := range distractors {
		ids = append(ids, record.ID)
	}

	// Shuffle the IDs and determine the new index of the correct answer
	shuffledIDs, newCorrectIdx := shuffleIDs(ids, correctRecord.ID)

	return &domain.DefinitionQuizQuestion{
		Definition: correctRecord.Definition,
		IDs:        shuffledIDs,
		CorrectIdx: newCorrectIdx,
	}
}

// generateWriteDefinitionQuizQuestion generates a quiz question where you need to write the definition
func (qs *Service) generateWriteDefinitionQuizQuestion() *domain.WriteDefinitionQuizQuestion {
	if len(qs.records) == 0 {
		log.Fatal("No records available to generate a quiz question")
	}

	// Randomly select an ID
	correctIdx := rand.Intn(len(qs.records))
	correctRecord := qs.records[correctIdx]

	return &domain.WriteDefinitionQuizQuestion{
		ID: correctRecord.ID,
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

// shuffleIDs shuffles the IDs and returns the shuffled IDs along with the new index of the correct ID
func shuffleIDs(ids []string, correctID string) ([]string, int) {
	shuffledIDs := make([]string, len(ids))
	copy(shuffledIDs, ids)
	rand.Shuffle(len(shuffledIDs), func(i, j int) {
		shuffledIDs[i], shuffledIDs[j] = shuffledIDs[j], shuffledIDs[i]
	})

	var newCorrectIdx int
	for i, id := range shuffledIDs {
		if id == correctID {
			newCorrectIdx = i
			break
		}
	}

	return shuffledIDs, newCorrectIdx
}
