package artist

import (
	"github.com/labstack/echo/v5"
)

func Routes(e *echo.Echo, h *Handler) {
	e.GET("/artist/:id", h.GetArtist)
}
