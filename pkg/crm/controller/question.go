package crm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
	questionStore *QuestionStore
}

func NewQuestionController(questionStore *QuestionStore) *QuestionController {
	return &QuestionController{
		questionStore: questionStore,
	}
}

func (controller *QuestionController) GetQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	questions, err := controller.questionStore.GetQuestions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
}

func (controller *QuestionController) GetQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

func (controller *QuestionController) CreateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	var question Question
	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdQuestion, err := controller.questionStore.CreateQuestion(question.ID, question.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdQuestion)
}

func (controller *QuestionController) UpdateQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

func (controller *QuestionController) DeleteQuestionHandler(w http.ResponseWriter, r *http.Request) {
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

func (store *QuestionStore) GetQuestions() ([]*Question, error) {
	rows, err := store.db.Query("SELECT id, content FROM questions")
	if err != nil {
		return nil, fmt.Errorf("could not get questions: %v", err)
	}
	defer rows.Close()

	var questions []*Question
	for rows.Next() {
		var question Question
		err := rows.Scan(&question.ID, &question.Content)
		if err != nil {
			return nil, fmt.Errorf("could not scan question: %v", err)
		}
		questions = append(questions, &question)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("could not iterate over questions: %v", err)
	}

	return questions, nil
}

func (store *QuestionStore) GetQuestion(id int) (*Question, error) {
	var question Question
	err := store.db.QueryRow("SELECT id, content FROM questions WHERE id = $1", id).Scan(&question.ID, &question.Content)
	if err != nil {
		return nil, fmt.Errorf("could not get question: %v", err)
	}
	return &question, nil
}

func (store *QuestionStore) CreateQuestion(id int, content string) (*Question, error) {
	err := store.db.QueryRow("INSERT INTO questions(content) VALUES($1) RETURNING ID", content).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("could not create question: %v", err)
	}
	return &Question{ID: id, Content: content}, nil
}

func (store *QuestionStore) UpdateQuestion(id int, content string) error {
	_, err := store.db.Exec("UPDATE questions SET content = $1 WHERE id = $2", content, id)
	if err != nil {
		return fmt.Errorf("could not update question: %v", err)
	}
	return nil
}

func (store *QuestionStore) DeleteQuestion(id int) error {
	_, err := store.db.Exec("DELETE FROM questions WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("could not delete question: %v", err)
	}
	return nil
}
