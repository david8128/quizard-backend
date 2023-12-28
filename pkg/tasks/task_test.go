package tasks

import (
	"testing"
)

func TestExecuteTask(t *testing.T) {
	task := Task{
		Name:   "Test task",
		Script: "/path/to/script.sh",
	}

	err := task.Run()
	if err != nil {
		t.Errorf("ExecuteTask() failed, expected %v, got %v", nil, err)
	}
}
