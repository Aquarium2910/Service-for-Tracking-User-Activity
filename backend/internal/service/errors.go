package service

import "errors"

var (
	ErrInvalidEvent  = errors.New("event cannot be nil")
	ErrInvalidUserID = errors.New("user_id must be greater than 0")
	ErrInvalidAction = errors.New("action cannot be empty")

	ErrInvalidFilter = errors.New("filter cannot be nil")
	ErrInvalidDates  = errors.New("start date cannot be after end date")
	ErrMissingDates  = errors.New("start and end dates are required")
)
