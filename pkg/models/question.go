package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
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
