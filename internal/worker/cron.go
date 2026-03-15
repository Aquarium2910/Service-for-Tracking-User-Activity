package worker

import (
	"context"
	"time"

	"test/internal/service"
)

type ActivityWorker struct {
	svc      service.ActivityService
	interval time.Duration
}

func NewActivityWorker(svc service.ActivityService, interval time.Duration) *ActivityWorker {
	return &ActivityWorker{
		svc:      svc,
		interval: interval,
	}
}

func (w *ActivityWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.runAggregation(ctx)

	for {
		select {
		case <-ticker.C:
			w.runAggregation(ctx)
		case <-ctx.Done():
			return
		}
	}
}

func (w *ActivityWorker) runAggregation(ctx context.Context) {
	end := time.Now().Truncate(time.Hour)
	start := end.Add(-w.interval)

	err := w.svc.ProcessActivityStats(ctx, start, end)
	if err != nil {
		return
	}
}
