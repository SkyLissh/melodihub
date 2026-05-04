package search

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeSearchClient struct {
	result []deezer.Track
	err    error
	calls  int
	query  string
	limit  int
}

func (f *fakeSearchClient) Find(
	_ context.Context,
	query string,
	limit int,
) ([]deezer.Track, error) {
	f.calls++
	f.query = query
	f.limit = limit

	if f.err != nil {
		return nil, f.err
	}

	return f.result, nil
}

func TestServiceGetSearchBuildsResult(t *testing.T) {
	fakeClient := &fakeSearchClient{
		result: []deezer.Track{
			{
				ID:    101,
				Title: "Hysteria",
				Album: deezer.Album{
					ID:    201,
					Title: "Absolution",
					Cover: "https://example.com/album.jpg",
				},
				Artist: deezer.Artist{
					ID:      301,
					Name:    "Muse",
					Link:    "https://example.com/muse",
					Picture: "https://example.com/artist.jpg",
				},
			},
		},
	}

	service := NewService(fakeClient, nil)

	result, err := service.GetSearch(context.Background(), "muse", 5)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, 1, fakeClient.calls)
	assert.Equal(t, "muse", fakeClient.query)
	assert.Equal(t, 5, fakeClient.limit)

	require.Len(t, result.Tracks, 1)
	assert.Equal(t, "deezer:101", result.Tracks[0].ID)

	require.Len(t, result.Artists, 1)
	assert.Equal(t, "Muse", result.Artists[0].Name)

	require.Len(t, result.Albums, 1)
	assert.Equal(t, "Absolution", result.Albums[0].Title)

	require.NotNil(t, result.TopResult)
	assert.NotEmpty(t, result.TopResult.ID)
}

func TestHandlerGetSearchReturnsBadRequestWhenQueryMissing(t *testing.T) {
	handler := NewHandler(NewService(&fakeSearchClient{}, nil))
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/search", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetSearch(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
	assert.Equal(t, http.StatusBadRequest, apiErr.Code)
	assert.Equal(t, "Query parameter 'q' is required", apiErr.Message)
}

func TestHandlerGetSearchUsesLimitAndReturnsJSON(t *testing.T) {
	fakeClient := &fakeSearchClient{
		result: []deezer.Track{
			{
				ID:    101,
				Title: "Hysteria",
				Album: deezer.Album{
					ID:    201,
					Title: "Absolution",
					Cover: "https://example.com/album.jpg",
				},
				Artist: deezer.Artist{
					ID:      301,
					Name:    "Muse",
					Link:    "https://example.com/muse",
					Picture: "https://example.com/artist.jpg",
				},
			},
		},
	}

	handler := NewHandler(NewService(fakeClient, nil))
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/search?q=muse&limit=3", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetSearch(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 1, fakeClient.calls)
	assert.Equal(t, 3, fakeClient.limit)

	var result Result
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	require.Len(t, result.Tracks, 1)
	assert.Equal(t, "Hysteria", result.Tracks[0].Title)
	require.NotNil(t, result.TopResult)
}
