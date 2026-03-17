// Package handlers contains the handler functions for the API routes
package handlers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/labstack/echo/v5"

	"github.com/skylissh/melodihub/internal/cache"
	"github.com/skylissh/melodihub/internal/middleware"
	"github.com/skylissh/melodihub/internal/model"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
)

type CarouselHandler interface {
	GetCarousels(c *echo.Context) error
}

type carouselHandler struct{}

func NewCarouselHandler() CarouselHandler {
	return &carouselHandler{}
}

// Get Carousels godoc
//
//	@Summary		Get carousels
//	@Description	Get all home carousels, to show initial content
//	@Tags			carousels
//	@Accept			json
//	@Produce		json
//	@Param			country	query		string	false	"Country"
//	@Success		200		{object}	model.Carousel
//	@Failure		500		{object}	model.APIError
//	@Router			/carousels [get]
func (h *carouselHandler) GetCarousels(c *echo.Context) error {
	ctx := c.Request().Context()

	p, err := middleware.GetProviders(c)
	if err != nil {
		return err
	}

	cacheClient, err := middleware.GetCache(c)
	if err != nil {
		return err
	}

	country := c.QueryParam("country")
	if country == "" {
		country = "united states"
	}

	cacheKey := fmt.Sprintf("carousels:%s", country)
	if cached, err := cache.GetJSON[[]model.Carousel](
		ctx, cacheClient, cacheKey,
	); err == nil {
		return c.JSON(http.StatusOK, cached)
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		c.Logger().Error("Cache get failed", "error", err)
	}

	g1, _ := errgroup.WithContext(c.Request().Context())

	lastfmTopArtists := make([]lastfm.Artist, 0, 10)
	lastfmTopTracks := make([]lastfm.Track, 0, 10)

	g1.Go(func() error {
		result, err := p.Lastfm.Geo.GetTopArtists(country, 10, 1)
		if err != nil {
			return err
		}

		lastfmTopArtists = result
		return nil
	})

	g1.Go(func() error {
		result, err := p.Lastfm.Geo.GetTopTracks(country, 10, 1)
		if err != nil {
			return err
		}

		lastfmTopTracks = result
		return nil
	})

	if err := g1.Wait(); err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Message: "Error fetching Last.fm data",
			Code:    http.StatusInternalServerError,
		})
	}

	g, _ := errgroup.WithContext(context.Background())
	g.SetLimit(5)

	deezerArtists := make([]deezer.SearchResult, len(lastfmTopArtists))
	deezerTracks := make([]deezer.SearchResult, len(lastfmTopTracks))

	for i, artist := range lastfmTopArtists {
		g.Go(func() error {
			i, artist := i, artist // capture loop variables

			deezerArtist, err := p.Deezer.Search.Find(artist.Name, 10)
			if err != nil {
				return err
			}
			deezerArtists[i] = deezerArtist[0]
			return nil
		})
	}

	for i, track := range lastfmTopTracks {
		g.Go(func() error {
			deezerTrack, err := p.Deezer.Search.Find(track.Name, 10)
			if err != nil {
				return err
			}

			deezerTracks[i] = deezerTrack[0]
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Message: "Error fetching Deezer data",
			Code:    http.StatusInternalServerError,
		})
	}

	artists := make([]model.CarouselItem, len(deezerArtists))
	for i, result := range deezerArtists {
		artists[i] = model.CarouselItem{
			ID:   fmt.Sprintf("deezer:%d", result.Artist.ID),
			Type: model.ArtistType,
			Name: result.Artist.Name,
			Images: []model.Image{
				{
					URL: result.Artist.Picture,
				},
			},
		}
	}

	tracks := make([]model.CarouselItem, len(deezerTracks))
	for i, result := range deezerTracks {
		tracks[i] = model.CarouselItem{
			ID:   fmt.Sprintf("deezer:%d", result.ID),
			Type: model.AlbumType,
			Name: result.Title,
			Images: []model.Image{
				{
					URL: result.Album.Cover,
				},
			},
		}
	}

	carousels := []model.Carousel{
		{
			Title: "Top Artists",
			Items: artists,
		},
		{
			Title: "Top Tracks",
			Items: tracks,
		},
	}

	err = cache.SetJSON(ctx, cacheClient, cacheKey, carousels, cache.DefaultTTL)
	if err != nil {
		c.Logger().Error("Cache set failed", "error", err)
	}

	return c.JSON(http.StatusOK, carousels)
}
