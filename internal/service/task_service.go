package service

import (
	"context"
	"errors"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/repository"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidTask  = errors.New("invalid task")
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

func (service *TaskService) GetTasks(ctx context.Context, userID uint16) ([]entity.Task, error) {
	// ctx := context.Background()
	return service.repo.FindAll(ctx, userID)
}

func (service *TaskService) GetSingleTask(ctx context.Context, id, userID uint16) (entity.Task, error) {
	return service.repo.FindById(ctx, id, userID)
	// if err != nil {
	// 	return task, err
	// }
	// return task, err
}

func (service *TaskService) AddTask(ctx context.Context, task entity.Task) (entity.Task, error) {
	// ctx := context.Background()

	return service.repo.Insert(ctx, task)
	// if err != nil {
	// 	return err
	// }
	// return nil
}

func (service *TaskService) UpdateTask(ctx context.Context, id, userID uint16, task entity.Task) (entity.Task, error) {
	task, err := service.repo.Update(ctx, task, id, userID)
	if err != nil {
		return task, err
	}

	return task, nil
}

func (service *TaskService) DeleteTask(ctx context.Context, id, userID uint16) error {
	return service.repo.Delete(ctx, id, userID)
}
