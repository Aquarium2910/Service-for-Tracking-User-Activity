package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"test/internal/database"
	"test/internal/models"
)

var (
	ErrInvalidEvent  = errors.New("event cannot be nil")
	ErrInvalidUserID = errors.New("user_id must be greater than 0")
	ErrInvalidAction = errors.New("action cannot be empty")

	ErrInvalidFilter = errors.New("filter cannot be nil")
	ErrInvalidDates  = errors.New("start date cannot be after end date")
	ErrMissingDates  = errors.New("start and end dates are required")
)

type ActivityService interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEvents(ctx context.Context, filter *models.EventFilter) ([]models.Event, error)
	ProcessActivityStats(ctx context.Context, start time.Time, end time.Time) error
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
	if event == nil {
		return fmt.Errorf("activityService.CreateEvent - %w", ErrInvalidEvent)
	}

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

func (s *activityService) GetEvents(ctx context.Context, filter *models.EventFilter) ([]models.Event, error) {
	if filter == nil {
		return nil, fmt.Errorf("activityService.GetEvents - filter error: %w", ErrInvalidFilter)
	}

	if filter.UserID <= 0 {
		return nil, fmt.Errorf("activityService.GetEvents - filter error: %w", ErrInvalidUserID)
	}

	if !filter.StartDate.IsZero() && !filter.EndDate.IsZero() {
		if filter.StartDate.After(filter.EndDate) {
			return nil, fmt.Errorf("activityService.GetEvents - filter error: %w", ErrInvalidDates)
		}
	}

	events, err := s.repo.GetEvents(ctx, filter)

	if err != nil {
		return nil, fmt.Errorf("activityService.GetEvents - repository error: %w", err)
	}

	return events, nil
}

func (s *activityService) ProcessActivityStats(ctx context.Context, start time.Time, end time.Time) error {
	if start.IsZero() || end.IsZero() {
		return fmt.Errorf("activityService.ProcessActivityStats - validation error: %w", ErrMissingDates)
	}

	if start.After(end) {
		return fmt.Errorf("activityService.ProcessActivityStats - validation error: %w", ErrInvalidDates)
	}

	err := s.repo.AggregateActivity(ctx, start, end)
	if err != nil {
		return fmt.Errorf("activityService.ProcessActivityStats - repository error: %w", err)
	}

	return nil
}
