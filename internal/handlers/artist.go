package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/samber/lo"
	"golang.org/x/sync/errgroup"

	"github.com/skylissh/melodihub/internal/cache"
	"github.com/skylissh/melodihub/internal/mapper"
	"github.com/skylissh/melodihub/internal/middleware"
	"github.com/skylissh/melodihub/internal/model"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
)

type ArtistHandler interface {
	GetArtist(c *echo.Context) error
}

type artistHandler struct{}

func NewArtistHandler() ArtistHandler {
	return &artistHandler{}
}

// Get Artist Details godoc
//
//	@Summary		Get artist details
//	@Description	Get detailed information about an artist, including top tracks, albums, and similar artists
//	@Tags			artist
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"Artist ID (e.g., deezer:123)"
//	@Success		200	{object}	model.Artist
//	@Failure		400	{object}	model.APIError
//	@Failure		500	{object}	model.APIError
//	@Router			/artist/{id} [get]
func (h artistHandler) GetArtist(c *echo.Context) error {
	ctx := c.Request().Context()

	// Deezer ID example: deezer:123, we need to extract the numeric part for API calls
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Message: "Artist ID is required",
			Code:    http.StatusBadRequest,
		})
	}

	dID, err := url.PathUnescape(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Message: "Invalid artist ID format",
			Code:    http.StatusBadRequest,
		})
	}
	id = dID

	p, err := middleware.GetProviders(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Message: "Failed to get providers",
			Code:    http.StatusInternalServerError,
		})
	}

	cacheClient, err := middleware.GetCache(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Message: "Failed to get cache client",
			Code:    http.StatusInternalServerError,
		})
	}

	cacheKey := fmt.Sprintf("artist:%s", id)
	if cached, err := cache.GetJSON[model.Artist](
		ctx, cacheClient, cacheKey,
	); err == nil {
		return c.JSON(http.StatusOK, cached)
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		c.Logger().Error("Cache get failed", "error", err)
	}

	deezerArtist, err := p.Deezer.Artist.GetByID(ctx, strings.TrimPrefix(id, "deezer:"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Message: "Failed to fetch artist from Deezer",
			Code:    http.StatusInternalServerError,
		})
	}

	// Use errgroup to fetch Last.fm artist details, top tracks, and top albums concurrently
	g1, g1Ctx := errgroup.WithContext(ctx)
	g1.SetLimit(5)

	lastfmArtist := &lastfm.ArtistDetail{}
	lastfmTopTracks := make([]lastfm.Track, 0, 10)
	lastfmTopAlbums := make([]lastfm.Album, 0, 15)

	g1.Go(func() error {
		detail, err := p.Lastfm.Artist.GetDetail(g1Ctx, deezerArtist.Name)
		if err != nil {
			return err
		}

		*lastfmArtist = *detail
		return nil
	})

	g1.Go(func() error {
		tracks, err := p.Lastfm.Artist.GetTopTracks(g1Ctx, deezerArtist.Name, 10)
		if err != nil {
			return err
		}

		lastfmTopTracks = tracks.Tracks
		return nil
	})

	g1.Go(func() error {
		albums, err := p.Lastfm.Artist.GetTopAlbums(g1Ctx, deezerArtist.Name, 15)

		if err != nil {
			return err
		}

		lastfmTopAlbums = albums.Album
		return nil
	})

	if err := g1.Wait(); err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Message: "Failed to fetch artist details from Last.fm",
			Code:    http.StatusInternalServerError,
		})
	}

	g2, g2Ctx := errgroup.WithContext(ctx)
	g2.SetLimit(5)

	similar := make([]model.SimpleArtist, len(lastfmArtist.Similar.Artists))
	topTracks := make([]model.SimpleTrack, 0, 10)
	topAlbums := make([]model.SimpleAlbum, 0, 15)

	// Fetch Deezer search results for the artist to find matching tracks and albums
	g2.Go(func() error {
		result, err := p.Deezer.Search.Find(g2Ctx, deezerArtist.Name, 50)
		if err != nil {
			return err
		}

		search := mapper.SearchFromDeezer(result, nil)

		for _, a := range search.Albums {
			if len(topAlbums) >= 15 {
				break
			}

			_, ok := lo.Find(lastfmTopAlbums, func(la lastfm.Album) bool {
				return strings.EqualFold(la.Name, a.Title)
			})
			if !ok {
				continue
			}

			topAlbums = append(topAlbums, model.SimpleAlbum{
				ID:     a.ID,
				Title:  a.Title,
				Artist: model.SimpleArtist{Name: a.Artist.Name, ID: a.Artist.ID},
				Images: []model.Image{{URL: a.Image}},
			})
		}

		for _, t := range search.Tracks {
			if len(topTracks) >= 10 {
				break
			}

			m, ok := lo.Find(lastfmTopTracks, func(lt lastfm.Track) bool {
				return strings.EqualFold(lt.Name, t.Title)
			})
			if !ok {
				continue
			}

			topTracks = append(topTracks, model.SimpleTrack{
				ID:        t.ID,
				Title:     t.Title,
				Duration:  t.Duration,
				Artist:    model.SimpleArtist{Name: t.Artist.Name, ID: t.Artist.ID},
				Images:    []model.Image{{URL: t.Image}},
				Rank:      t.Rank,
				Listeners: m.Listeners,
			})
		}

		return nil
	})

	for i, a := range lastfmArtist.Similar.Artists {
		g2.Go(func() error {
			i, a := i, a // capture loop variables

			result, err := p.Deezer.Search.Find(g2Ctx, a.Name, 1)
			if err != nil {
				return err
			}

			if len(result) == 0 {
				return nil // Skip if no results found
			}

			artist := result[0].Artist

			similar[i] = model.SimpleArtist{
				ID:     fmt.Sprintf("deezer:%d", artist.ID),
				Name:   artist.Name,
				Images: []model.Image{{URL: artist.Picture}},
			}

			return nil
		})
	}

	if err := g2.Wait(); err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Message: "Failed to fetch similar artists from Deezer",
			Code:    http.StatusInternalServerError,
		})
	}

	artist := model.Artist{
		ID:        fmt.Sprintf("deezer:%d", deezerArtist.ID),
		Name:      deezerArtist.Name,
		Listeners: lastfmArtist.Stats.Listeners,
		Bio:       lastfmArtist.Bio.Summary,
		Images:    []model.Image{{URL: deezerArtist.Picture}},
		TopTracks: topTracks,
		TopAlbums: topAlbums,
		Similar:   similar,
	}

	err = cache.SetJSON(ctx, cacheClient, cacheKey, artist, cache.DefaultTTL)
	if err != nil {
		c.Logger().Error("Cache set failed", "error", err)
	}

	return c.JSON(http.StatusOK, artist)
}
