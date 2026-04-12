package artist

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/labstack/echo/v5"

	common "github.com/skylissh/melodihub/internal/domain"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

// Get Artist Details godoc
//
//	@Summary		Get artist details
//	@Description	Get detailed information about an artist, including top tracks, albums, and similar artists
//	@Tags			artist
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Artist ID (e.g., deezer:123)"
//	@Success		200	{object}	Detail
//	@Failure		400	{object}	common.APIError
//	@Failure		500	{object}	common.APIError
//	@Router			/artist/{id} [get]
func (h *Handler) GetArtist(c *echo.Context) error {
	ctx := c.Request().Context()

	// Deezer ID example: deezer:123, we need to extract the numeric part for API calls
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, common.NewAPIError(
			"Artist ID is required",
			http.StatusBadRequest,
		))
	}

	id, err := url.PathUnescape(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, common.NewAPIError(
			"Invalid artist ID format",
			http.StatusBadRequest,
		))
	}

	artist, err := h.service.GetArtistDetails(ctx, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, common.NewAPIError(
			fmt.Sprintf("Failed to get artist details: %v", err),
			http.StatusInternalServerError,
		))
	}

	return c.JSON(http.StatusOK, artist)
}
