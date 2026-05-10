package search

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

// Get Search Results godoc
//
//	@Summary		Get search results
//	@Description	Get search results based on a query string
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Param			q		query		string	true	"Search query"
//	@Param			limit	query		int		false	"Number of results to return"
//	@Success		200		{object}	ResultResponse
//	@Failure		400		{object}	domain.APIError
//	@Failure		500		{object}	domain.APIError
//	@Router			/search [get]
func (h *Handler) GetSearch(c *echo.Context) error {
	ctx := c.Request().Context()

	query, err := ParseQuery(c.QueryParam("q"))
	if err != nil {
		apiErr, status := domain.APIErrorFromMelodiError(domain.InvalidParam(err))
		return c.JSON(status, apiErr)
	}

	limit, err := ParseLimit(c.QueryParam("limit"))
	if err != nil {
		apiErr, status := domain.APIErrorFromMelodiError(domain.InvalidParam(err))
		return c.JSON(status, apiErr)
	}

	response, err := h.service.GetSearch(ctx, query, limit)
	if err != nil {
		apiErr, status := domain.APIErrorFromMelodiError(err)
		return c.JSON(status, apiErr)
	}

	return c.JSON(http.StatusOK, response)
}
