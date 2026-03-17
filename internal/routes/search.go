package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/skylissh/melodihub/internal/handlers"
)

func RegisterSearchRoutes(e *echo.Echo, h handlers.SearchHandler) {
	e.GET("/search", h.GetSearch)
	e.GET("/search/suggestions", h.GetSuggestions)
}
