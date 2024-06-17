package main

import (
	v1 "lsfb/internal/api/v1"
	"lsfb/internal/application/quizz"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	quizzService := quizz.NewQuizzService("output.csv")
	quizzHttpHandler := v1.NewQuizzHTTPHandler(quizzService)

	mux.HandleFunc("GET /api/quizz", quizzHttpHandler.HandleGetQuizz)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
