package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"test/internal/models"
	"test/internal/service"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	activityService ActivityService
	logger          *slog.Logger
}

func NewHandler(activityService ActivityService, logger *slog.Logger) *Handler {
	return &Handler{
		activityService: activityService,
		logger:          logger.With("component", "http_handler"),
	}
}

func (h *Handler) HandleCreateEvent(c echo.Context) error {
	var event models.Event

	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid JSON format",
		})
	}

	err := h.activityService.CreateEvent(c.Request().Context(), &event)

	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidAction) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		h.logger.Error("unexpected error in HandleCreateEvent", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error while saving event",
		})
	}

	return c.JSON(http.StatusCreated, event)
}

func (h *Handler) HandleGetEvent(c echo.Context) error {
	var filter models.EventFilter

	if err := c.Bind(&filter); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid query parameters",
		})
	}

	events, err := h.activityService.GetEvents(c.Request().Context(), &filter)

	if err != nil {
		if errors.Is(err, service.ErrInvalidUserID) || errors.Is(err, service.ErrInvalidFilter) ||
			errors.Is(err, service.ErrInvalidDates) {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}

		h.logger.Error("unexpected error in HandleGetEvent", slog.Any("error", err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error while getting event",
		})
	}

	return c.JSON(http.StatusOK, events)
}
