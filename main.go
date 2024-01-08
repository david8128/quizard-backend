package main

import (
	"net/http"

	"github.com/david8128/quizard-backend/pkg/controllers"
	"github.com/david8128/quizard-backend/pkg/tasks"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()

	// Rutas para el CRM
	r.HandleFunc("/questions", controllers.GetQuestions).Methods("GET")
	r.HandleFunc("/questions/{id}", controllers.GetQuestion).Methods("GET")
	r.HandleFunc("/questions", controllers.CreateQuestion).Methods("POST")
	r.HandleFunc("/questions/{id}", controllers.UpdateQuestion).Methods("PUT")
	r.HandleFunc("/questions/{id}", controllers.DeleteQuestion).Methods("DELETE")

	// Rutas para las tareas de validación
	r.HandleFunc("/check_config", tasks.CheckConfig).Methods("POST")

	http.ListenAndServe(":8000", r)
}
