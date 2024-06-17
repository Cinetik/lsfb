package v1

import (
	"encoding/json"
	"lsfb/internal/domain"
	"net/http"
)

type QuizzService interface {
	GetQuizz() domain.QuizInterface
	//CheckAnswer(quizz Quizz) bool
}

type QuizzError struct {
	StatusCode int    `json:"status_code"`
	Error      string `json:"error"`
}

type QuizzHTTPHandler struct {
	quizzService QuizzService
}

func NewQuizzHTTPHandler(quizzService QuizzService) *QuizzHTTPHandler {
	return &QuizzHTTPHandler{quizzService: quizzService}
}

func (handler *QuizzHTTPHandler) HandleGetQuizz(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	quiz := handler.quizzService.GetQuizz()

	if err := json.NewEncoder(w).Encode(quiz); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
