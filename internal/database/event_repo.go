package database

import (
	"context"
	"test/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepo struct {
	db *pgxpool.Pool
}

func NewEventRepo(db *pgxpool.Pool) *EventRepo {
	return &EventRepo{db: db}
}

func (r *EventRepo) Create(ctx context.Context, event *models.Event) error {
	query := `
		INSERT INTO events (user_id, action, metadata)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := r.db.QueryRow(ctx, query, event.UserID, event.Action, event.Metadata).Scan(&event.ID, &event.CreatedAt)

	return err
}
