package storage

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
	// "github.com/WatShitTooYaa/go-task-manager-api/intenal/entity"
)

func (storage *Storage) HandleList() func([]string) error {
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
			if task.Completed {
				completed = "X"
			} else {
				completed = " "
			}
			result := fmt.Sprintf("%d. [%s] %s (%s) \t %s", task.Id, completed, task.Content, task.Priority, task.Timestamp)
			fmt.Println(result)
		}

		return nil
	}
}

func (storage *Storage) HandleAdd() func([]string) error {
	return func(cmds []string) error {
		tasks, err := storage.Load()
		if err != nil {
			return err
		}

		var lastId uint16 = 0
		if len(tasks) > 0 {
			lastId = tasks[len(tasks)-1].Id
		}

		task := entity.Task{
			Id:        lastId + 1,
			Content:   cmds[1],
			Completed: false,
			Timestamp: time.Now().Format(time.RFC3339),
			// Time:        time.Now().Format("02-01-2006 15:04:05"),
			Priority: "low",
		}

		tasks = append(tasks, task)

		return storage.Save(tasks)
	}
}

func (storage *Storage) HandleUpdate() func([]string) error {
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
					tasks[i].Completed = isCompleted
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

func (storage *Storage) HandleDelete() func([]string) error {
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

		var tempTask []entity.Task
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
