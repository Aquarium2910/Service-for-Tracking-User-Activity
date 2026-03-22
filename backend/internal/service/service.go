package service

import (
	"context"
	"fmt"
	"time"

	"test/internal/models"
)

type activityService struct {
	repo EventRepo
}

func NewActivityService(repo EventRepo) *activityService {
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
