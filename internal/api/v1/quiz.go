package v1

import (
	"encoding/json"
	"lsfb/internal/domain"
	"net/http"
)

type QuizService interface {
	GetQuiz() domain.QuizQuestionType
	//CheckAnswer(quiz Quizz) bool
}

type QuizError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type QuizHTTPHandler struct {
	quizService QuizService
}

func NewQuizHTTPHandler(quizService QuizService) *QuizHTTPHandler {
	return &QuizHTTPHandler{quizService: quizService}
}

func (handler *QuizHTTPHandler) HandleGetQuiz(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	quiz := handler.quizService.GetQuiz()
	wrapped := domain.WrapQuizQuestion(quiz)

	if err := json.NewEncoder(w).Encode(wrapped); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
