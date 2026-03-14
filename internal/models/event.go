package models

import (
	"encoding/json"
	"time"
)

type Event struct {
	ID        int64           `json:"id"`
	UserID    int             `json:"user_id"`
	Action    string          `json:"action"`
	Metadata  json.RawMessage `json:"metadata"`
	CreatedAt time.Time       `json:"created_at"`
}

type EventFilter struct {
	UserID    int       `query:"user_id"`
	StartDate time.Time `query:"start_date"`
	EndDate   time.Time `query:"end_date"`
}
