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
	limit  uint
}

func (f *fakeSearchClient) Find(
	_ context.Context,
	query string,
	limit uint,
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

	response, err := service.GetSearch(context.Background(), "muse", 5)

	require.NoError(t, err)
	require.NotNil(t, response)

	assert.Equal(t, 1, fakeClient.calls)
	assert.Equal(t, "muse", fakeClient.query)
	assert.Equal(t, uint(5), fakeClient.limit)

	require.Len(t, response.Tracks, 1)
	assert.Equal(t, "deezer:101", response.Tracks[0].ID)

	require.Len(t, response.Artists, 1)
	assert.Equal(t, "Muse", response.Artists[0].Name)

	require.Len(t, response.Albums, 1)
	assert.Equal(t, "Absolution", response.Albums[0].Title)

	assert.NotEmpty(t, response.TopResult.ID)
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
	assert.Equal(t, domain.CodeInvalidParam, apiErr.Code)
	assert.Equal(t, "Invalid param", apiErr.Message)
	assert.Equal(t, "search query cannot be empty", apiErr.Cause)
}

func TestHandlerGetSearchReturnsBadRequestWhenLimitInvalid(t *testing.T) {
	handler := NewHandler(NewService(&fakeSearchClient{}, nil))
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/search?q=muse&limit=0", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetSearch(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
	assert.Equal(t, domain.CodeInvalidParam, apiErr.Code)
	assert.Equal(t, "Invalid param", apiErr.Message)
	assert.Equal(t, "search limit is too small (min 1)", apiErr.Cause)
}

func TestHandlerGetSearchMapsDeezerInvalidQueryToInvalidParam(t *testing.T) {
	apiErr, status := getSearchErrorResponse(t, deezer.DeezerError{
		Kind: deezer.ErrQueryInvalid,
		Msg:  "query must contain at least one token",
	})

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, domain.CodeInvalidParam, apiErr.Code)
	assert.Equal(t, "Invalid param", apiErr.Message)
	assert.Equal(t, "query invalid: query must contain at least one token", apiErr.Cause)
}

func TestHandlerGetSearchMapsDeezerNotFoundToNotFound(t *testing.T) {
	apiErr, status := getSearchErrorResponse(t, deezer.DeezerError{Kind: deezer.ErrNotFound})

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, domain.CodeNotFound, apiErr.Code)
	assert.Equal(t, "Not found", apiErr.Message)
	assert.Equal(t, "not found", apiErr.Cause)
}

func TestHandlerGetSearchMapsDeezerRateLimitedErrors(t *testing.T) {
	tests := []error{
		deezer.DeezerError{Kind: deezer.ErrQuotaExceeded},
		deezer.DeezerError{Kind: deezer.ErrItemsLimit},
	}

	for _, providerErr := range tests {
		t.Run(providerErr.Error(), func(t *testing.T) {
			apiErr, status := getSearchErrorResponse(t, providerErr)

			assert.Equal(t, http.StatusTooManyRequests, status)
			assert.Equal(t, domain.CodeRateLimited, apiErr.Code)
			assert.Equal(t, "Rate limited", apiErr.Message)
			assert.Equal(t, "deezer: "+providerErr.Error(), apiErr.Cause)
		})
	}
}

func TestHandlerGetSearchMapsDeezerProviderUnavailableErrors(t *testing.T) {
	tests := []error{
		deezer.DeezerError{Kind: deezer.ErrServiceBusy},
		deezer.DeezerError{Kind: deezer.ErrInvalidResponse},
		deezer.DeezerError{Kind: deezer.ErrServer},
		deezer.DeezerError{Kind: deezer.ErrUnknownCode},
	}

	for _, providerErr := range tests {
		t.Run(providerErr.Error(), func(t *testing.T) {
			apiErr, status := getSearchErrorResponse(t, providerErr)

			assert.Equal(t, http.StatusBadGateway, status)
			assert.Equal(t, domain.CodeProviderUnavailable, apiErr.Code)
			assert.Equal(t, "Provider unavailable", apiErr.Message)
			assert.Equal(t, "deezer: "+providerErr.Error(), apiErr.Cause)
		})
	}
}

func getSearchErrorResponse(t *testing.T, providerErr error) (domain.APIError, int) {
	t.Helper()

	handler := NewHandler(NewService(&fakeSearchClient{err: providerErr}, nil))
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/search?q=muse", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetSearch(c)
	require.NoError(t, err)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))

	return apiErr, rec.Code
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
	assert.Equal(t, uint(3), fakeClient.limit)

	var result struct {
		TopResult *TopResultResponse      `json:"top_result"`
		Tracks    []TrackSummaryResponse  `json:"tracks"`
		Artists   []ArtistSummaryResponse `json:"artists"`
		Albums    []AlbumSummaryResponse  `json:"albums"`
	}
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &result))
	require.Len(t, result.Tracks, 1)
	assert.Equal(t, "Hysteria", result.Tracks[0].Title)
	require.NotNil(t, result.TopResult)
}

func TestResponseFromResultReturnsErrorForNilResult(t *testing.T) {
	response, err := ResponseFromResult(nil)

	require.Error(t, err)
	assert.Nil(t, response)
}

func TestResultResponseJSONRoundTripPreservesSearchData(t *testing.T) {
	result, err := ResponseFromResult(ResultFromDeezer("muse", []deezer.Track{
		{
			ID:       101,
			Title:    "Hysteria",
			Duration: 227,
			Album: deezer.Album{
				ID:    201,
				Title: "Absolution",
				Cover: "https://example.com/album.jpg",
			},
			Artist: deezer.Artist{
				ID:      301,
				Name:    "Muse",
				Picture: "https://example.com/artist.jpg",
			},
		},
	}))
	require.NoError(t, err)

	payload, err := json.Marshal(result)
	require.NoError(t, err)

	var restored ResultResponse
	require.NoError(t, json.Unmarshal(payload, &restored))

	require.Len(t, restored.Tracks, 1)
	assert.Equal(t, "deezer:101", restored.Tracks[0].ID)
	assert.Equal(t, "Muse", restored.Tracks[0].Artists[0].Name)

	require.Len(t, restored.Artists, 1)
	assert.Equal(t, "Muse", restored.Artists[0].Name)

	require.Len(t, restored.Albums, 1)
	assert.Equal(t, "Absolution", restored.Albums[0].Title)
	assert.Equal(t, "Muse", restored.Albums[0].Artists[0].Name)
}
