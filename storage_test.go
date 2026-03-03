package main

import (
	"os"
	"testing"
)

func TestAddTask(t *testing.T) {
	var fileTestName = "add_test"
	storage := NewStorage(fileTestName)
	defer os.Remove(fileTestName)

	tasks := [][]string{
		{`add "kentang"`, "kentang"},
		{`add "ngejim"`, "ngejim"},
		{`add "batur"`, "batur"},
	}

	for i, task := range tasks {
		// var expected error =nil
		handler := storage.handleAdd()

		result := handler(task)

		if result != nil {
			t.Errorf("task-%d, result = %s, want nil", i+1, result.Error())
		}
	}
}

func TestUpdateTask(t *testing.T) {
	var fileTestName = "update_test"
	storage := NewStorage(fileTestName)
	defer os.Remove(fileTestName)

	initTask := [][]string{
		{`add "kentang"`, "tidur"},
		{`add "ngejim"`, "bangun"},
		{`add "batur"`, "makan"},
		{`add "batur"`, "minum"},
		{`add "batur"`, "solat"},
	}

	for i, task := range initTask {
		// var expected error =nil
		handler := storage.handleAdd()

		result := handler(task)

		if result != nil {
			t.Errorf("task-%d, result = %s, want nil", i+1, result.Error())
		}
	}

	tasks := [][]string{
		{`update 5 -t "walk"`, `1`, `t`, `"walk"`},
		{`update 4 -p "medium"`, `2`, `p`, `"medium"`},
		{`update 3 -s "true"`, `3`, `s`, `"true"`},
		{`update 29 -t "true"`, `29`, `t`, `"true"`},
		{`update 3 -p "true"`, `5`, `p`, `"true"`},
	}

	for i, task := range tasks {
		// var expected error =nil
		handler := storage.handleUpdate()

		result := handler(task)

		if result != nil {
			t.Errorf("task-%d, result = %s, want nil", i+1, result.Error())
		}
	}
}

func TestDelete(t *testing.T) {
	var fileTestName = "delete_test"
	storage := NewStorage(fileTestName)
	defer os.Remove(fileTestName)

	initTask := [][]string{
		{`add "kentang"`, "tidur"},
		{`add "ngejim"`, "bangun"},
		{`add "batur"`, "makan"},
	}

	for i, task := range initTask {
		// var expected error =nil
		handler := storage.handleAdd()

		result := handler(task)

		if result != nil {
			t.Errorf("task-%d, result = %s, want nil", i+1, result.Error())
		}
	}

	deleteTasks := [][]string{
		{`delete 2`, `1`},
		{`delete 5`, `2`},
		{`delete 2998`, `2998`},
	}

	for i, task := range deleteTasks {
		// var expected error =nil
		handler := storage.handleDelete()

		result := handler(task)

		if result != nil {
			t.Errorf("task-%d, result = %s, want nil", i+1, result.Error())
		}
	}
}
