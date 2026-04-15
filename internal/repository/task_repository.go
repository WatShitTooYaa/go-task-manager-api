package repository

import (
	"context"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
)

type TaskRepository interface {
	Insert(ctx context.Context, task entity.Task) (entity.Task, error)
	FindAll(ctx context.Context, userID uint16) ([]entity.Task, error)
	FindById(ctx context.Context, id, userID uint16) (entity.Task, error)
	Update(ctx context.Context, newTask entity.Task, id, userID uint16) (entity.Task, error)
	Delete(ctx context.Context, id, userID uint16) error
}
