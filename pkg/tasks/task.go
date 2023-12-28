package tasks

import (
	"os/exec"
	"log"
)

type Task struct {
	Name string
	Script string
}

func NewTask(name string, script string) *Task {
	return &Task{
		Name: name,
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