package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
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

func (storage *Storage) handleList() func([]string) error {
	return func(cmds []string) error {
		tasks, err := storage.Load()
		if err != nil {
			return err
		}

		if len(tasks) == 0 {
			fmt.Println("Tidak ada task yang tersedia")
			return nil
		}

		for _, task := range tasks {
			completed := ""
			if task.IsCompleted {
				completed = "X"
			} else {
				completed = " "
			}
			result := fmt.Sprintf("%d. [%s] %s (%s) \t %s", task.Id, completed, task.Content, task.Priority, task.Time)
			fmt.Println(result)
		}

		return nil
	}
}

func (storage *Storage) handleAdd() func([]string) error {
	return func(cmds []string) error {
		tasks, err := storage.Load()
		if err != nil {
			return err
		}

		var lastId uint16 = 0
		if len(tasks) > 0 {
			lastId = tasks[len(tasks)-1].Id
		}

		task := Task{
			Id:          lastId + 1,
			Content:     cmds[1],
			IsCompleted: false,
			Time:        time.Now().Format(time.RFC3339),
			// Time:        time.Now().Format("02-01-2006 15:04:05"),
			Priority: "low",
		}

		tasks = append(tasks, task)

		return storage.Save(tasks)
	}
}

func (storage *Storage) handleUpdate() func([]string) error {
	return func(cmds []string) error {
		tasks, err := storage.Load()
		if err != nil {
			return err
		}

		rawId := cmds[1]
		id, err := strconv.Atoi(rawId)
		if err != nil {
			return err
		}
		flag := cmds[2]
		val := cmds[3]

		var found = false
		for i, task := range tasks {
			if task.Id == uint16(id) {
				found = true
				switch flag {
				case "t", "task":
					tasks[i].Content = val
				case "p", "priority":
					if val == "low" || val == "medium" || val == "high" {
						tasks[i].Priority = val
					} else {
						return errors.New("priority salah")
					}
				case "s", "status":
					isCompleted, err := strconv.ParseBool(val)
					if err != nil {
						return err
					}
					tasks[i].IsCompleted = isCompleted
					// break loopTask
				default:
					return errors.New("flag tidak ditemukan")
				}
				break
			}
		}
		if !found {
			return errors.New("id tidak ditemukan")
		}

		return storage.Save(tasks)
	}
}

func (storage *Storage) handleDelete() func([]string) error {
	return func(cmds []string) error {
		rawId := cmds[1]
		id, err := strconv.Atoi(rawId)
		if err != nil {
			return err
		}

		tasks, err := storage.Load()
		if err != nil {
			return err
		}

		var tempTask []Task
		var found = false
		for _, task := range tasks {
			if task.Id == uint16(id) {
				found = true
				continue
			}
			tempTask = append(tempTask, task)
		}
		if !found {
			return errors.New("id not found")
		}

		return storage.Save(tempTask)
	}
}
