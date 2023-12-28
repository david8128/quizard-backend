package main

import (
	"net/http"

	"mi-backend/pkg/crm"
	"mi-backend/pkg/tasks"

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
	r.HandleFunc("/tasks", tasks.GetTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", tasks.GetTask).Methods("GET")
	r.HandleFunc("/tasks", tasks.CreateTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", tasks.UpdateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", tasks.DeleteTask).Methods("DELETE")

	http.ListenAndServe(":8000", r)
}
