package service

import (
	"context"

	"github.com/WatShitTooYaa/go-task-manager-api/internal/entity"
	"github.com/WatShitTooYaa/go-task-manager-api/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) LoginService(ctx context.Context, user entity.UserParam) (entity.User, error) {
	return s.repo.Login(ctx, user)
}

func (s *UserService) RegisterService(ctx context.Context, user entity.UserParam) (entity.User, error) {
	return s.repo.Insert(ctx, user)
}

// func ()  {

// }
