package tasks

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type Task struct {
	Name   string
	Script string
}

func NewTask(name string, script string) *Task {
	return &Task{
		Name:   name,
		Script: script,
	}
}

func (t *Task) Run() error {
	cmd := exec.Command("/bin/bash", t.Script)

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running task %s: %v", t.Name, err)
		return err
	}

	log.Printf("Task %s completed successfully", t.Name)
	return nil
}

func CheckConfig(w http.ResponseWriter, r *http.Request) {
	// Add your configuration checking logic here
	// Return an error if the configuration is invalid

	// Check if the script exists for the given task ID
	taskID := r.FormValue("taskID")
	scriptPath := fmt.Sprintf("/path/to/scripts/%s.sh", taskID)
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		http.Error(w, "Not valid conf", http.StatusBadRequest)
		return
	}

	// Run the validation task by executing the bash script
	cmd := exec.Command("/bin/bash", scriptPath)
	output, err := cmd.Output()
	if err != nil {
		http.Error(w, "Not valid conf", http.StatusBadRequest)
		return
	}

	// Check if the output contains the word "success"
	if strings.Contains(string(output), "success") {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Accepted"))
	} else {
		http.Error(w, "Not valid conf", http.StatusBadRequest)
	}
}
