package handlers

import (
	"context"
	"test/internal/models"
)

type ActivityService interface {
	CreateEvent(ctx context.Context, event *models.Event) error
	GetEvents(ctx context.Context, filter *models.EventFilter) ([]models.Event, error)
}
