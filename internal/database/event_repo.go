package database

import (
	"context"
	"fmt"
	"test/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EventRepo interface {
	Create(ctx context.Context, event *models.Event) error
	GetEvents(ctx context.Context, filter *models.EventFilter) ([]models.Event, error)
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

func (r *eventRepo) GetEvents(ctx context.Context, filter *models.EventFilter) ([]models.Event, error) {
	query := `SELECT id, user_id, action, metadata, created_at FROM events WHERE user_id = $1`

	args := []any{filter.UserID}

	argId := 2

	if !filter.StartDate.IsZero() {
		query += fmt.Sprintf(" AND created_at >= $%d", argId)
		args = append(args, filter.StartDate)
		argId++
	}
	if !filter.EndDate.IsZero() {
		query += fmt.Sprintf(" AND created_at <= $%d", argId)
		args = append(args, filter.EndDate)
		argId++
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("EventRepo.GetEvents - failed to execute query: %w", err)
	}
	defer rows.Close()

	var events []models.Event
	for rows.Next() {
		var e models.Event
		err := rows.Scan(&e.ID, &e.UserID, &e.Action, &e.Metadata, &e.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("EventRepo.GetEvents - failed to scan row: %w", err)
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("EventRepo.GetEvents - rows iteration error: %w", err)
	}

	return events, nil
}
