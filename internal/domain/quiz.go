package domain

type QuizInterface interface {
	Type() string
}

// QuizQuestion represents a quiz question with signs
type QuizQuestion struct {
	ID         string   `json:"id"`
	Signs      []string `json:"signs"`
	CorrectIdx int      `json:"correct_idx"`
}

// DefinitionQuizQuestion represents a quiz question with definitions
type DefinitionQuizQuestion struct {
	Definition string   `json:"definition"`
	IDs        []string `json:"ids"`
	CorrectIdx int      `json:"correct_idx"`
}

// WriteDefinitionQuizQuestion represents a quiz question where you need to write the definition
type WriteDefinitionQuizQuestion struct {
	ID string `json:"id"`
}

func (q *QuizQuestion) Type() string {
	return "SignsQuiz"
}

func (dq *DefinitionQuizQuestion) Type() string {
	return "DefinitionQuiz"
}

func (wd *WriteDefinitionQuizQuestion) Type() string {
	return "WriteDefinitionQuiz"
}
