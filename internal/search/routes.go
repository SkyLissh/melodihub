package search

import (
	"github.com/labstack/echo/v5"
)

func Routes(e *echo.Echo, h *Handler) {
	e.GET("/search", h.GetSearch)
}
