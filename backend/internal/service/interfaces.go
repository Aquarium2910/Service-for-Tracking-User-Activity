package service

import (
	"context"
	"test/internal/models"
	"time"
)

type EventRepo interface {
	Create(ctx context.Context, event *models.Event) error
	GetEvents(ctx context.Context, filter *models.EventFilter) ([]models.Event, error)
	AggregateActivity(ctx context.Context, start time.Time, end time.Time) error
}
