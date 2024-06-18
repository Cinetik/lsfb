package domain

// QuizQuestionType Interface for quiz question types
type QuizQuestionType interface {
	GetType() string
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

func (q *QuizQuestion) GetType() string {
	return "SignsQuiz"
}

func (dq *DefinitionQuizQuestion) GetType() string {
	return "DefinitionQuiz"
}

func (wd *WriteDefinitionQuizQuestion) GetType() string {
	return "WriteDefinitionQuiz"
}

// QuizQuestionWrapper Wrapper struct for JSON encoding
type QuizQuestionWrapper struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// WrapQuizQuestion wraps a quiz question with its type
func WrapQuizQuestion(q QuizQuestionType) QuizQuestionWrapper {
	return QuizQuestionWrapper{
		Type: q.GetType(),
		Data: q,
	}
}
