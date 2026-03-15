package main

import (
	"context"
	"log"
	"os"
	"test/internal/worker"
	"time"

	"test/internal/api"
	"test/internal/database"
	"test/internal/service"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	_ = godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	}

	ctx := context.Background()

	dbPool, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("database is unreachable: %v", err)
	}

	repo := database.NewEventRepo(dbPool)
	svc := service.NewActivityService(repo)
	handler := api.NewHTTPHandler(svc)

	activityWorker := worker.NewActivityWorker(svc, 4*time.Hour)
	go activityWorker.Start(ctx)

	e := echo.New()

	v1 := e.Group("/api/v1")
	v1.POST("/events", handler.HandleCreateEvent)
	v1.GET("/events", handler.HandleGetEvent)

	log.Println("Server is running on port :8080")
	if err := e.Start(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
