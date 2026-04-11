package repository

import (
	"context"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
)

type UserRepository interface {
	Insert(ctx context.Context, user entity.UserParam) (entity.User, error)
	Login(ctx context.Context, user entity.UserParam) (entity.User, error)
	Get(ctx context.Context, id uint16) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
	Update(ctx context.Context, newUser entity.UserParam, id uint16) (entity.User, error)
	Delete(ctx context.Context, id uint16) error
}
