package storage

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
	// "github.com/WatShitTooYaa/go-task-manager-api/internal/storage"
)

type Storage struct {
	fileName string
}

type TaskJob struct {
	Content  string
	Priority string
}

type TaskResult struct {
	Task  entity.Task
	Error error
}

func NewStorage(filename string) *Storage {
	return &Storage{
		fileName: filename,
	}
}

func (storage *Storage) Load() ([]entity.Task, error) {
	file, err := os.ReadFile(storage.fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []entity.Task{}, nil
		}
		return nil, err
	}

	task := []entity.Task{}
	if len(file) > 0 {

		err = json.Unmarshal(file, &task)
		if err != nil {
			return nil, err
		}
	}

	return task, nil
}

func (storage *Storage) Save(tasks []entity.Task) error {
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

func (storage *Storage) GetByID(id uint16) (*entity.Task, error) {
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

func (s *Storage) AddTask(content, priority string) (*entity.Task, error) {
	tasks, err := s.Load()
	if err != nil {
		return nil, err
	}

	newID := uint16(len(tasks) + 1)
	task := entity.Task{
		Id:        newID,
		Content:   content,
		Priority:  priority,
		Completed: false,
		Timestamp: time.Now().Format(time.RFC3339),
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
			tasks[i].Completed = completed
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

	newTasks := []entity.Task{}
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

func worker(id int, jobs chan TaskJob, result chan TaskResult) {
	// var mutex sync.Mutex

	for job := range jobs {

		if job.Content == "" {
			result <- TaskResult{
				Task:  entity.Task{},
				Error: errors.New("Task kosong"),
			}
			continue
		}

		if job.Priority != "low" && job.Priority != "medium" && job.Priority != "high" {

			result <- TaskResult{
				Task:  entity.Task{},
				Error: errors.New("Priority tidak sesuai"),
			}
			continue
		}
		task := entity.Task{
			Content:   job.Content,
			Priority:  job.Priority,
			Completed: false,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		result <- TaskResult{
			Task:  task,
			Error: nil,
		}
	}
}

func (storage *Storage) ImportCSV() func([]string) error {
	return func(cmds []string) error {

		filename := strings.TrimSpace(cmds[1])
		fmt.Println("filename : ", filename)

		file, err := os.Open(filename)
		if err != nil {
			return errors.New("File tidak ditemukan")
		}
		defer file.Close()

		reader := csv.NewReader(file)
		// reader.FieldsPerRecord = -1
		// reader.TrimLeadingSpace = true

		workerAmount := 5
		// defer close(chanJobs)
		// defer close(chanResult)

		records, err := reader.ReadAll()
		if err != nil {
			return err
		}
		numJobs := len(records) - 1
		chanJobs := make(chan TaskJob, numJobs)
		chanResult := make(chan TaskResult, numJobs)

		//run worker
		for i := 1; i <= workerAmount; i++ {
			go worker(i, chanJobs, chanResult)
		}

		//mengirim
		// i := 0
		for i, record := range records {
			if i == 0 {
				continue
			}
			chanJobs <- TaskJob{
				Content:  record[0],
				Priority: record[1],
			}
		}
		close(chanJobs)

		// mutex := sync.Mutex{}

		validTask := make([]entity.Task, 0, numJobs)

		//menerima result
		// lastId := len(mainTasks)
		// i = 0

		numErrors := 0

		for range numJobs {
			result := <-chanResult
			if result.Error != nil {
				numErrors++
				fmt.Printf("task : %s error\n", result.Error)

				// return nil
				continue
			}

			// task := result.Task
			// task.Id = uint16(j)

			fmt.Printf("task : %s success\n", result.Task.Content)
			validTask = append(validTask, result.Task)

		}

		close(chanResult)

		mainTasks, err := storage.Load()
		if err != nil {
			return err
		}

		validId := len(mainTasks) + 1
		for j := range validTask {
			// fmt.Printf("id : %d success\n", validId +1)

			validTask[j].Id = uint16(validId + j)
		}

		mainTasks = append(mainTasks, validTask...)
		err = storage.Save(mainTasks)
		if err != nil {
			return err
		}

		fmt.Printf("\nImport complete:\n")
		fmt.Printf("  Valid: %d tasks\n", len(validTask))
		fmt.Printf("  Errors: %d tasks\n", numErrors)
		return nil
	}
}
