package database

import (
	"context"
	"fmt"
	"test/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepo interface {
	Create(ctx context.Context, event *models.Event) error
}

type eventRepo struct {
	db *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool) EventRepo {
	return &eventRepo{db: db}
}

func (r *eventRepo) Create(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO events (user_id, action, metadata)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, event.UserID, event.Action, event.Metadata).Scan(&event.ID, &event.CreatedAt)
	if err != nil {
		return fmt.Errorf("EventRepo.Create - failed to insert event: %w", err)
	}

	return nil
}
