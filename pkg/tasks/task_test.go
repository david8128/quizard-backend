package tasks

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExecuteTask(t *testing.T) {
	task := Task{
		Name:   "Test task",
		TaskID: "1",
	}

	err := task.Run()
	if err != nil {
		t.Errorf("ExecuteTask() failed, expected %v, got %v", nil, err)
	}

	// Test CheckConfig
	reqBody := strings.NewReader(`{"TaskId": "1"}`)
	req, err := http.NewRequest("POST", "/config", reqBody)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CheckConfig)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("CheckConfig() returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := "Accepted"
	if rr.Body.String() != expected {
		t.Errorf("CheckConfig() returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}
