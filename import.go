package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

type TaskJob struct {
	Content  string
	Priority string
}

type TaskResult struct {
	Task  Task
	Error error
}

func worker(id int, jobs chan TaskJob, result chan TaskResult) {
	// var mutex sync.Mutex

	for job := range jobs {

		if job.Content == "" {
			result <- TaskResult{
				Task:  Task{},
				Error: errors.New("Task kosong"),
			}
			continue
		}

		if job.Priority != "low" && job.Priority != "medium" && job.Priority != "high" {

			result <- TaskResult{
				Task:  Task{},
				Error: errors.New("Priority tidak sesuai"),
			}
			continue
		}
		task := Task{
			Content:     job.Content,
			Priority:    job.Priority,
			IsCompleted: false,
			Time:        time.Now().Format(time.RFC3339),
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

		validTask := make([]Task, 0, numJobs)

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
