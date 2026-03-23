package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v5"
	"github.com/skylissh/melodihub/internal/cache"
	"github.com/skylissh/melodihub/internal/mapper"
	"github.com/skylissh/melodihub/internal/middleware"
	"github.com/skylissh/melodihub/internal/model"
	"github.com/skylissh/melodihub/internal/utils"
)

type SearchHandler interface {
	GetSearch(c *echo.Context) error
	GetSuggestions(c *echo.Context) error
}

type searchHandler struct{}

func NewSearchHandler() SearchHandler {
	return &searchHandler{}
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
//	@Success		200		{object}	model.SearchResult
//	@Failure		400		{object}	model.APIError
//	@Failure		500		{object}	model.APIError
//	@Router			/search [get]
func (h searchHandler) GetSearch(c *echo.Context) error {
	ctx := c.Request().Context()

	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Message: "Query parameter 'q' is required",
			Code:    http.StatusBadRequest,
		})
	}

	limit := 10
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

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

	cacheKey := fmt.Sprintf("search:%s:limit:%d", strings.ToLower(query), limit)
	if cached, err := cache.GetJSON[model.SearchResult](
		ctx, cacheClient, cacheKey,
	); err == nil {
		return c.JSON(http.StatusOK, cached)
	}

	if !errors.Is(err, cache.ErrCacheMiss) {
		c.Logger().Error("Cache get failed", "error", err)
	}

	results, err := p.Deezer.Search.Find(ctx, query, limit)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Message: "Failed to search for suggestions",
			Code:    http.StatusInternalServerError,
		})
	}

	search := mapper.SearchFromDeezer(results, &query)

	err = cache.SetJSON(ctx, cacheClient, cacheKey, search, cache.DefaultTTL)
	if err != nil {
		c.Logger().Error("Cache set failed", "error", err)
	}

	return c.JSON(http.StatusOK, search)
}

// Get Search Suggestions godoc
//
//	@Summary		Get search suggestions
//	@Description	Get search suggestions based on a query string
//	@Tags			search
//	@Accept			json
//	@Produce		json
//	@Param			q		query		string	true	"Search query"
//	@Param			limit	query		int		false	"Number of suggestions to return"
//	@Success		200		{array}		string
//	@Failure		400		{object}	model.APIError
//	@Failure		500		{object}	model.APIError
//	@Router			/search/suggestions [get]
func (h searchHandler) GetSuggestions(c *echo.Context) error {
	ctx := c.Request().Context()

	query := c.QueryParam("q")
	if query == "" {
		return c.JSON(http.StatusBadRequest, model.APIError{
			Message: "Query parameter 'q' is required",
			Code:    http.StatusBadRequest,
		})
	}

	limit := 5
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			limit = l
		}
	}

	p, err := middleware.GetProviders(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Message: "Failed to get providers",
			Code:    http.StatusInternalServerError,
		})
	}

	results, err := p.Deezer.Search.Find(ctx, query, 10)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, model.APIError{
			Message: "Failed to search for suggestions",
			Code:    http.StatusInternalServerError,
		})
	}

	suggest := make(map[string]struct{})

	for _, result := range results {
		track := fmt.Sprintf("%s - %s", result.Title, result.Artist.Name)
		album := fmt.Sprintf("%s - %s", result.Album.Title, result.Artist.Name)

		if _, exists := suggest[track]; !exists {
			suggest[strings.ToLower(track)] = struct{}{}
		}

		if _, exists := suggest[result.Artist.Name]; !exists {
			suggest[strings.ToLower(result.Artist.Name)] = struct{}{}
		}

		if _, exists := suggest[album]; !exists {
			suggest[strings.ToLower(album)] = struct{}{}
		}
	}

	suggestions := make([]string, 0, limit)
	for s := range suggest {
		if len(suggestions) >= cap(suggestions) {
			break
		}

		suggestions = append(suggestions, s)
	}

	slices.SortFunc(suggestions, func(a, b string) int {
		aMatch := utils.RankString(query, a)
		bMatch := utils.RankString(query, b)

		if aMatch > bMatch {
			return -1
		} else if aMatch < bMatch {
			return 1
		}

		return 0
	})

	return c.JSON(http.StatusOK, suggestions)
}
