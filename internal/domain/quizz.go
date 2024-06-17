package domain

type QuizInterface interface {
	GetID() string
	GetType() string
}

type Quiz struct {
	ID   string
	Type string
}

func (q Quiz) GetID() string {
	return q.ID
}

func (q Quiz) GetType() string {
	return q.Type
}

type SelectionQuiz struct {
	Quiz
	GifID   string
	Signe   []string
	Correct string
}

type SigneQuiz struct {
	Quiz
	GifID   string
	Signe   []string
	Correct string
}

type WriteQuiz struct {
	Quiz
	GifID string
}
