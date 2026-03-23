package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/skylissh/melodihub/internal/handlers"
)

func RegisterArtistRoutes(e *echo.Echo, h handlers.ArtistHandler) {
	e.GET("/artist/:id", h.GetArtist)
}
