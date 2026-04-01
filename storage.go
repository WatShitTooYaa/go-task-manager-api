package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Storage struct {
	fileName string
}

func NewStorage(filename string) *Storage {
	return &Storage{
		fileName: filename,
	}
}

func (storage *Storage) Load() ([]Task, error) {
	file, err := os.ReadFile(storage.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []Task{}, nil
		}
		return nil, err
	}

	task := []Task{}
	if len(file) > 0 {

		err = json.Unmarshal(file, &task)
		if err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (storage *Storage) Save(tasks []Task) error {
	json, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		return err
	}
	err = os.WriteFile(storage.fileName, json, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (storage *Storage) GetByID(id uint16) (*Task, error) {
	tasks, err := storage.Load()
	if err != nil {
		return nil, err
	}

	// found := false
	for _, task := range tasks {
		if task.Id == id {
			return &task, nil
		}
	}

	err = fmt.Errorf("task with ID = %d not found\n", id)

	return nil, err
}

func (s *Storage) AddTask(content, priority string) (*Task, error) {
	tasks, err := s.Load()
	if err != nil {
		return nil, err
	}

	newID := uint16(len(tasks) + 1)
	task := Task{
		Id:          newID,
		Content:     content,
		Priority:    priority,
		IsCompleted: false,
		Time:        time.Now().Format(time.RFC3339),
	}

	tasks = append(tasks, task)
	err = s.Save(tasks)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (s *Storage) UpdateTask(id uint16, content, priority string, completed bool) error {
	tasks, err := s.Load()
	if err != nil {
		return err
	}

	found := false
	for i, task := range tasks {
		if task.Id == id {
			tasks[i].Content = content
			tasks[i].Priority = priority
			tasks[i].IsCompleted = completed
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("task with id %d not found", id)
	}

	return s.Save(tasks)
}

func (s *Storage) DeleteTask(id uint16) error {
	tasks, err := s.Load()
	if err != nil {
		return err
	}

	newTasks := []Task{}
	found := false
	for _, task := range tasks {
		if task.Id == id {
			found = true
			continue
		}
		newTasks = append(newTasks, task)
	}

	if !found {
		return fmt.Errorf("task with id %d not found", id)
	}

	return s.Save(newTasks)
}
