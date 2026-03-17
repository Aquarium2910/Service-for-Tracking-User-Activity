package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	h.setupCORS(e)

	v1 := e.Group("/api/v1")
	v1.POST("/events", h.HandleCreateEvent)
	v1.GET("/events", h.HandleGetEvent)
}

func (h *Handler) setupCORS(e *echo.Echo) {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType},
	}))
}
