package worker

import (
	"context"
	"time"
)

type WorkerService interface {
	ProcessActivityStats(ctx context.Context, start time.Time, end time.Time) error
}
