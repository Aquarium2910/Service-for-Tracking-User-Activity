package service

import (
	"context"
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
		return ErrInvalidEvent
	}

	if event.UserID <= 0 {
		return ErrInvalidUserID
	}

	if event.Action == "" {
		return ErrInvalidAction
	}

	err := s.repo.Create(ctx, event)
	if err != nil {
		return err
	}

	return nil
}

func (s *activityService) GetEvents(ctx context.Context, filter *models.EventFilter) ([]models.Event, error) {
	if filter == nil {
		return nil, ErrInvalidFilter
	}

	if filter.UserID <= 0 {
		return nil, ErrInvalidUserID
	}

	if !filter.StartDate.IsZero() && !filter.EndDate.IsZero() {
		if filter.StartDate.After(filter.EndDate) {
			return nil, ErrInvalidDates
		}
	}

	events, err := s.repo.GetEvents(ctx, filter)

	if err != nil {
		return nil, err
	}

	return events, nil
}

func (s *activityService) ProcessActivityStats(ctx context.Context, start time.Time, end time.Time) error {
	if start.IsZero() || end.IsZero() {
		return ErrMissingDates
	}

	if start.After(end) {
		return ErrInvalidDates
	}

	err := s.repo.AggregateActivity(ctx, start, end)
	if err != nil {
		return err
	}

	return nil
}
