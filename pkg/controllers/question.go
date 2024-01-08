package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/david8128/quizard-backend/pkg/models"
	"github.com/gorilla/mux"
)

type Question struct {
	ID      int
	Content string
}

type QuestionStore struct {
	db *sql.DB
}

type QuestionController struct {
	questionStore *models.QuestionStore
}

func NewQuestionController() *QuestionController {
	questionStore := models.NewQuestionStore()

	return &QuestionController{
		questionStore: questionStore,
	}
}

func (controller *QuestionController) GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := controller.questionStore.GetQuestions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func (controller *QuestionController) GetQuestion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	questionID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	question, err := controller.questionStore.GetQuestion(questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}

func (controller *QuestionController) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var question Question
	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdQuestion, err := controller.questionStore.CreateQuestion(question.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdQuestion)
}

func (controller *QuestionController) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	content := r.FormValue("content") // Get the question content from the request

	questionID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = controller.questionStore.UpdateQuestion(questionID, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (controller *QuestionController) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	questionID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = controller.questionStore.DeleteQuestion(questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
