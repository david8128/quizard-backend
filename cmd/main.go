package main

import (
	"net/http"

	"github.com/david8128/quizard-backend/pkg/controllers"
	"github.com/david8128/quizard-backend/pkg/tasks"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	qc := controllers.NewQuestionController()

	// Rutas para el CRM
	r.HandleFunc("/questions", qc.GetQuestions).Methods("GET")
	r.HandleFunc("/questions/{id}", qc.GetQuestion).Methods("GET")
	r.HandleFunc("/questions", qc.CreateQuestion).Methods("POST")
	r.HandleFunc("/questions/{id}", qc.UpdateQuestion).Methods("PUT")
	r.HandleFunc("/questions/{id}", qc.DeleteQuestion).Methods("DELETE")

	// Rutas para las tareas de validación
	r.HandleFunc("/check_config", tasks.CheckConfig).Methods("POST")

	http.ListenAndServe(":8000", r)
}
