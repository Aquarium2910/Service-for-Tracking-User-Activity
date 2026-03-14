package service

import (
	"context"
	"errors"
	"fmt"
	"test/internal/database"
	"test/internal/models"
)

var (
	ErrInvalidUserID = errors.New("user_id must be greater than 0")
	ErrInvalidAction = errors.New("action cannot be empty")
)

type ActivityService interface {
	CreateEvent(ctx context.Context, event *models.Event) error
}

type activityService struct {
	repo database.EventRepo
}

func NewActivityService(repo database.EventRepo) ActivityService {
	return &activityService{
		repo: repo,
	}
}

func (s *activityService) CreateEvent(ctx context.Context, event *models.Event) error {
	if event.UserID <= 0 {
		return ErrInvalidUserID
	}

	if event.Action == "" {
		return ErrInvalidAction
	}

	err := s.repo.Create(ctx, event)
	if err != nil {
		return fmt.Errorf("activityService.CreateEvent - repository error: %w", err)
	}

	return nil
}
