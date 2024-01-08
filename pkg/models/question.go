package crm

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"strconv"

	db "github.com/david8128/quizard-backend/pkg/db"
)

type Question struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
}

type QuestionStore struct {
	db *sql.DB
}

func NewQuestionStore(database *sql.DB) *QuestionStore {

	db, err := sql.Open("postgres", "host=localhost port=5432 user=your_username password=your_password dbname=your_database_name sslmode=disable")
	if err != nil {
		error_msg := fmt.Sprintf("Failed to open database: %v", err)
		panic(error_msg)
	}

	err = db.Ping()
	if err != nil {
		error_msg := fmt.Sprintf("Failed to ping database: %v", err)
		panic(error_msg)
	}

	if db == nil {
		panic("db is nil")
	}
	return &QuestionStore{
		db: db,
	}
}

func (store *QuestionStore) Close() {
	store.db.Close()
}

func (store *QuestionStore) GetQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	questions, err := store.GetQuestions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(questions)
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

func (store *QuestionStore) CreateQuestion(content string) (*Question, error) {
	var id int
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

func GetQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	questionStore := NewQuestionStore(db.GetDB())
	defer questionStore.Close()
	questionStore.GetQuestionsHandler(w, r)
}

func GetQuestionHandler(w http.ResponseWriter, r *http.Request) {
	questionStore := NewQuestionStore(db.GetDB())
	defer questionStore.Close()
	params := mux.Vars(r)
	id := params["id"]

	questionID, err := strconv.Atoi(id)
	if err != nil {
		body := fmt.Sprintf("Invalid ID: %v, and id is %s, params: %v", err, id, params["id"])
		http.Error(w, body, http.StatusBadRequest)
		return
	}

	question, err := questionStore.GetQuestion(questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(question)
}

func CreateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	questionStore := NewQuestionStore(db.GetDB())
	defer questionStore.Close()

	var question Question
	err := json.NewDecoder(r.Body).Decode(&question)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdQuestion, err := questionStore.CreateQuestion(question.Content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdQuestion)
}

func UpdateQuestionHandler(w http.ResponseWriter, r *http.Request) {
	questionStore := NewQuestionStore(db.GetDB())
	defer questionStore.Close()
	params := mux.Vars(r)
	id := params["id"]
	content := r.FormValue("content") // Get the question content from the request

	questionID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	err = questionStore.UpdateQuestion(questionID, content)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteQuestionHandler(w http.ResponseWriter, r *http.Request) {
	questionStore := NewQuestionStore(db.GetDB())
	defer questionStore.Close()
	params := mux.Vars(r)
	id := params["id"]

	questionID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = questionStore.DeleteQuestion(questionID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}