package carousel

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v5"

	common "github.com/skylissh/melodihub/internal/domain"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

// Get Carousels godoc
//
//	@Summary		Get carousels
//	@Description	Get all home carousels, to show initial content
//	@Tags			carousels
//	@Accept			json
//	@Produce		json
//	@Param			country	query		string	false	"Country"
//	@Success		200		{array}		Section
//	@Failure		500		{object}	common.APIError
//	@Router			/carousels [get]
func (h *Handler) GetCarousels(c *echo.Context) error {
	ctx := c.Request().Context()

	country := c.QueryParam("country")
	if country == "" {
		country = "united states"
	}

	carousels, err := h.service.GetCarousels(ctx, country)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewAPIError(
			fmt.Sprintf("Failed to get carousels: %v", err),
			http.StatusInternalServerError,
		))
	}

	return c.JSON(http.StatusOK, carousels)
}
