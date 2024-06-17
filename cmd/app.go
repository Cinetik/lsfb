package main

import (
	v1 "lsfb/internal/api/v1"
	"lsfb/internal/application/quiz"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	quizzService := quiz.NewQuizService("output.csv")
	quizzHttpHandler := v1.NewQuizHTTPHandler(quizzService)

	mux.HandleFunc("GET /api/quizz", quizzHttpHandler.HandleGetQuiz)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
