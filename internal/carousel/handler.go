package carousel

import (
	"net/http"

	"github.com/labstack/echo/v5"

	"github.com/skylissh/melodihub/internal/domain"
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
//	@Success		200		{array}		CarouselResponse
//	@Failure		400		{object}	domain.APIError
//	@Failure		404		{object}	domain.APIError
//	@Failure		429		{object}	domain.APIError
//	@Failure		502		{object}	domain.APIError
//	@Failure		500		{object}	domain.APIError
//	@Router			/carousels [get]
func (h *Handler) GetCarousels(c *echo.Context) error {
	ctx := c.Request().Context()

	country, err := ParseCountry(c.QueryParam("country"))
	if err != nil {
		apiErr, status := domain.APIErrorFromMelodiError(domain.InvalidParam(err))
		return c.JSON(status, apiErr)
	}

	carousels, err := h.service.GetCarousels(ctx, country)
	if err != nil {
		apiErr, status := domain.APIErrorFromMelodiError(err)
		return c.JSON(status, apiErr)
	}

	return c.JSON(http.StatusOK, carousels)
}
