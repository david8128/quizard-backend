package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func InitializeDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=your_username password=your_password dbname=your_database_name sslmode=disable")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
func TestCreateQuestionAndGetQuestion(t *testing.T) {

	r := mux.NewRouter()

	// Rutas para el CRM
	r.HandleFunc("/questions/{id}", GetQuestion).Methods("GET")

	question := Question{
		ID:      1,
		Content: "What is your name?",
	}

	// Convert question to JSON
	jsonQuestion, err := json.Marshal(question)
	if err != nil {
		t.Fatal(err)
	}

	// Create a new HTTP request with the JSON payload
	req, err := http.NewRequest("POST", "/questions", bytes.NewBuffer(jsonQuestion))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateQuestion)

	// Serve the HTTP request to the recorder
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: response %s got %v want %v",
			rr.Body.String(), status, http.StatusOK)
	}

	// Check the response body
	expected := `{"message": "Question created successfully"}`

	req, err = http.NewRequest("GET", "/questions/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetQuestion)

	// Serve the HTTP request to the recorder
	r.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: response %s got %v want %v",
			rr.Body.String(), status, http.StatusOK)
	}

	// Check the response body

	expected = `{"ID":1,"Content":"What is your name?"}`
	t.Logf("%x\n", strings.TrimSuffix(rr.Body.String(), "\n"))
	t.Logf("%x\n", expected)
	if strings.TrimSuffix(rr.Body.String(), "\n") != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
