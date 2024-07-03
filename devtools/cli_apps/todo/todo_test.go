package todo_test

import (
	"os"
	"testing"
	"todo-cli-app"
)

func TestAddd(t *testing.T) {
	l := todo.List{}
	l.Add("task 1")

	if l[0].Task != "task 1" {
		t.Errorf("Expected task 1, got %v", l[0].Task)
	}
}
func TestComplete(t *testing.T) {
	l := todo.List{}
	l.Add("task 1")

	err := l.Complete(1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if !l[0].Done {
		t.Errorf("Expected task to be done")
	}
}

func TestDelete(t *testing.T) {
	l := todo.List{}
	tasks := []string{
		"task 1",
		"task 2",
		"task 3",
	}
	for _, v := range tasks {
		l.Add(v)
	}
	if l[0].Task != tasks[0] {
		t.Errorf("Expected %v, got %v", tasks[0], l[0].Task)
	}
	l.Delete(2)
	if len(l) != 2 {
		t.Errorf("Expected 2, got %v", len(l))
	}
	if l[1].Task != tasks[2] {
		t.Errorf("Expected %v, got %v", tasks[2], l[1].Task)
	}

}

func TestSaveGet(t *testing.T) {
	l1 := todo.List{}
	l2 := todo.List{}

	l1.Add("task 1")

	if l1[0].Task != "task 1" {
		t.Errorf("Expected task 1, got %v", l1[0].Task)
	}
	tf, err := os.CreateTemp("", "")
	if err != nil {
		t.Fatalf("Error creating temp file: %s", err)
	}
	defer os.Remove(tf.Name())
	if err := l1.Save(tf.Name()); err != nil {
		t.Errorf("Error saving file: %s", err)
	}
	if err := l2.Get(tf.Name()); err != nil {
		t.Errorf("Error getting file: %s", err)
	}
	if l1[0].Task != l2[0].Task {
		t.Errorf("Expected %v, got %v", l1[0].Task, l2[0].Task)
	}
}
