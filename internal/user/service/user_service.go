package service

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"grpc-microservices/internal/user/domain/entity"
	"grpc-microservices/internal/user/domain/repository"
	"time"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) (*entity.User, error) {
	user := &entity.User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id string) (*entity.User, error) {
	return s.repo.GetByID(ctx, id)
}
