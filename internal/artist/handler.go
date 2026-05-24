package artist

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo/v5"

	"github.com/skylissh/melodihub/internal/domain"
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
//	@Success		200	{object}	DetailResponse
//	@Failure		400	{object}	domain.APIError
//	@Failure		404	{object}	domain.APIError
//	@Failure		429	{object}	domain.APIError
//	@Failure		502	{object}	domain.APIError
//	@Failure		500	{object}	domain.APIError
//	@Router			/artist/{id} [get]
func (h *Handler) GetArtist(c *echo.Context) error {
	ctx := c.Request().Context()

	idRaw, err := url.PathUnescape(c.Param("id"))
	if err != nil {
		apiErr, status := domain.APIErrorFromMelodiError(domain.InvalidParam(err))
		return c.JSON(status, apiErr)
	}

	id, err := domain.ParseID(idRaw)
	if err != nil {
		apiErr, status := domain.APIErrorFromMelodiError(domain.InvalidParam(err))
		return c.JSON(status, apiErr)
	}

	artist, err := h.service.GetArtistDetails(ctx, id)
	if err != nil {
		apiErr, status := domain.APIErrorFromMelodiError(err)
		return c.JSON(status, apiErr)
	}

	return c.JSON(http.StatusOK, artist)
}
