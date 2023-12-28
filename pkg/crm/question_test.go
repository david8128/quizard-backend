package crm

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateQuestionAndGetQuestion(t *testing.T) {
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
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := `{"message": "Question created successfully"}`
	// Initialize the QuestionStore
	questionStore := &QuestionStore{}

	db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest")
	if err != nil {
		t.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		t.Fatalf("Failed to ping database: %v", err)
	}

	// Set the db field of questionStore
	questionStore.db = db

	req, err = http.NewRequest("GET", "/questions/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(GetQuestion)

	// Serve the HTTP request to the recorder
	handler.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected = `{"id": 1, "question": "What is your name?"}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
