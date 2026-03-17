package models

import "time"

type ActivityStat struct {
	ID         int64     `json:"id"`
	UserID     int       `json:"user_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	EventCount int       `json:"event_count"`
	CreatedAt  time.Time `json:"created_at"`
}
