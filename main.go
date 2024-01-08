package main

import (
	"net/http"

	crm "github.com/david8128/quizard-backend/pkg/controllers"
	"github.com/david8128/quizard-backend/pkg/tasks"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Rutas para el CRM
	r.HandleFunc("/questions", crm.GetQuestions).Methods("GET")
	r.HandleFunc("/questions/{id}", crm.GetQuestion).Methods("GET")
	r.HandleFunc("/questions", crm.CreateQuestion).Methods("POST")
	r.HandleFunc("/questions/{id}", crm.UpdateQuestion).Methods("PUT")
	r.HandleFunc("/questions/{id}", crm.DeleteQuestion).Methods("DELETE")

	// Rutas para las tareas de validación
	r.HandleFunc("/check_config", tasks.CheckConfig).Methods("POST")

	http.ListenAndServe(":8000", r)
}
