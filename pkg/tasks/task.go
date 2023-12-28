package tasks

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	Name   string
	TaskID string
}

func NewTask(name string, taskid string) *Task {
	return &Task{
		Name:   name,
		TaskID: taskid,
	}
}

func (t *Task) Run() error {

	scriptPath := fmt.Sprintf("%s/../../scripts/task%s-validation.sh", os.Getenv("PWD"), t.TaskID)
	cmd := exec.Command("/bin/bash", scriptPath)

	err := cmd.Run()
	if err != nil {
		log.Panicf("Error running task %s: %v", t.Name, scriptPath)
		return err
	}

	log.Printf("Task %s completed successfully", t.Name)
	return nil
}

func CheckConfig(w http.ResponseWriter, r *http.Request) {
	// Add your configuration checking logic here
	// Return an error if the configuration is invalid

	// Check if the script exists for the given task ID
	var data map[string]interface{}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		panic(err)
	}
	TaskID := data["TaskId"].(string)
	scriptPath := fmt.Sprintf("%s/../../scripts/task%s-validation.sh", os.Getenv("PWD"), TaskID)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		error_msg := fmt.Sprintf("Doesn't exists: %v with ID %v", scriptPath, data["TaskID"])
		http.Error(w, error_msg, http.StatusBadRequest)
		return
	}

	// Run the validation task by executing the bash script
	cmd := exec.Command("/bin/bash", scriptPath)
	output, err := cmd.Output()
	if err != nil {
		error_msg := fmt.Sprintf("Error running task: %v", scriptPath)
		http.Error(w, error_msg, http.StatusBadRequest)
		return
	}

	// Check if the output contains the word "success"
	if strings.Contains(string(output), "success") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Accepted"))
	} else {
		http.Error(w, "Doesn't contains success", http.StatusBadRequest)
	}
}
