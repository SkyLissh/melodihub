package search

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v5"
	common "github.com/skylissh/melodihub/internal/domain"
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
//	@Success		200		{object}	Result
//	@Failure		400		{object}	common.APIError
//	@Failure		500		{object}	common.APIError
//	@Router			/search [get]
func (h *Handler) GetSearch(c *echo.Context) error {
	ctx := c.Request().Context()

	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, common.NewAPIError(
			"Query parameter 'q' is required",
			http.StatusBadRequest,
		))
	}

	limit := 10
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	search, err := h.service.GetSearch(ctx, query, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewAPIError(
			fmt.Sprintf("Failed to get search results: %v", err),
			http.StatusInternalServerError,
		))
	}

	return c.JSON(http.StatusOK, search)
}
