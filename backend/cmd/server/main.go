package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"test/internal/config"
	"test/internal/database"
	"test/internal/handlers"
	"test/internal/service"
	"test/internal/worker"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	cfg := config.LoadConfig(logger)

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("database is unreachable: %v", err)
	}

	repo := database.NewEventRepo(dbPool)
	svc := service.NewActivityService(repo)
	handler := handlers.NewHandler(svc, logger)

	e := echo.New()
	handler.RegisterRoutes(e)

	activityWorker := worker.NewActivityWorker(svc, cfg.WorkerInterval, logger)
	go activityWorker.Start(ctx)

	logger.Info("Server is running", slog.String("port", ":8080"))
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
