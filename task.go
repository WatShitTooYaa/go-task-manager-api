package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strconv"
)

type Task struct {
	Id uint16 `json:"id"`
	// Title       string `json:"title"`
	Content     string `json:"content"`
	IsCompleted bool   `json:"isCompleted"`
	Time        string `json:"timestamp"`
	Priority    string `json:"priority"`
}

type TaskJob interface {
	CreateTask()
	ReadTask()
	UpdateTask()
	DeleteTask()
}

const fileName = "storage.json"

func AddTask(task *Task) error {
	var tasks []Task

	bytes, err := os.ReadFile(fileName)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			tasks = []Task{}
		} else {
			return err
		}
		// panic(err.Error())
	}

	var lastId uint16 = 0
	if len(bytes) == 0 {
		tasks = []Task{}
	} else {
		// fmt.Println("masuk sini")
		err = json.Unmarshal(bytes, &tasks)
		if err != nil {
			// panic(err.Error())
			return err
		}
		lastId = tasks[len(tasks)-1].Id
	}

	task.Id = lastId + 1

	tasks = append(tasks, *task)

	output, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, output, 0644)
	if err != nil {
		return err

	}
	return nil
}

func ShowTask() error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var tasks []Task

	err = json.Unmarshal(file, &tasks)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		temp := fmt.Sprintf(`
		===== Task ====
		id : %v,
		task: %v,
		status : %v
		priority: %v,
		created at: %v,
		`, task.Id, task.Content, task.IsCompleted, task.Priority, task.Time)
		fmt.Println(temp)
	}
	return nil
}

func UpdateTask(id uint16, flag, newVal string) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	tasks := []Task{}
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		return err
	}
	// fmt.Println("len task update: ", len(tasks))

	// var tempErr error
	var idFound bool = false
loopTask:
	for i, task := range tasks {
		fmt.Printf("task id : %d, input id : %d\n", task.Id, id)
		if task.Id == id {
			idFound = true
			// fmt.Println("id ditemukan")
			switch flag {
			case "task":
				tasks[i].Content = newVal
				break loopTask
			case "priority":
				if newVal == "low" || newVal == "medium" || newVal == "high" {
					tasks[i].Priority = newVal
					break loopTask
				} else {
					return errors.New("priority salah")
				}
			case "status":
				isCompleted, err := strconv.ParseBool(newVal)
				if err != nil {
					return err
				}
				tasks[i].IsCompleted = isCompleted
				break loopTask
			default:
				return errors.New("flag tidak ditemukan")
			}
		} else {
			idFound = false
		}
	}
	if !idFound {
		return errors.New("id tidak ditemukan")
	}

	jsonRes, err := json.MarshalIndent(&tasks, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, jsonRes, os.FileMode(0644))
	if err != nil {
		return err
	}

	return nil
}

func DeleteTask(id uint16) error {
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	var tasks []Task
	var taskFound = false
	err = json.Unmarshal(file, &tasks)
	if err != nil {
		return err

	}
	newTask := []Task{}
	for _, task := range tasks {
		if task.Id == id {
			taskFound = true
			continue
		}
		newTask = append(newTask, task)
	}
	if taskFound {
		return errors.New("id tidak ditemukan")
	}

	res, err := json.MarshalIndent(newTask, "", "\t")
	if err != nil {
		return err
	}

	err = os.WriteFile(fileName, res, os.FileMode(0644))
	if err != nil {
		return err
	}
	return nil
}
