package routes

import (
	"github.com/labstack/echo/v5"
	"github.com/skylissh/melodihub/internal/handlers"
)

func RegisterCarouselsRoutes(e *echo.Echo, h handlers.CarouselHandler) {
	e.GET("/carousels", h.GetCarousels)
}
