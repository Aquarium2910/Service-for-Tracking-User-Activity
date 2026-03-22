package worker

import (
	"context"
	"log/slog"
	"time"
)

type ActivityWorker struct {
	svc      WorkerService
	interval time.Duration
	logger   *slog.Logger
}

func NewActivityWorker(svc WorkerService, interval time.Duration, logger *slog.Logger) *ActivityWorker {
	return &ActivityWorker{
		svc:      svc,
		interval: interval,
		logger:   logger.With("component", "worker"),
	}
}

func (w *ActivityWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	w.logger.Info("ActivityWorker started", slog.Duration("interval", w.interval))
	w.runAggregation(ctx)

	for {
		select {
		case <-ticker.C:
			w.runAggregation(ctx)
		case <-ctx.Done():
			w.logger.Info("ActivityWorker stopped")
			return
		}
	}
}

func (w *ActivityWorker) runAggregation(ctx context.Context) {
	end := time.Now().Truncate(time.Hour)
	start := end.Add(-w.interval)

	w.logger.Debug("running aggregation",
		slog.Time("start", start),
		slog.Time("end", end),
	)

	err := w.svc.ProcessActivityStats(ctx, start, end)
	if err != nil {
		w.logger.Error("aggregation failed", slog.Any("error", err))
		return
	}

	w.logger.Info("aggregation completed successfully")
}
