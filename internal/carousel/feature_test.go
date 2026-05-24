package carousel

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeGeoClient struct {
	artists   []lastfm.Artist
	tracks    []lastfm.Track
	artistErr error
	trackErr  error
	country   string
}

func (f *fakeGeoClient) GetTopArtists(
	_ context.Context,
	country string,
	_ int,
	_ int,
) ([]lastfm.Artist, error) {
	f.country = country
	if f.artistErr != nil {
		return nil, f.artistErr
	}
	return f.artists, nil
}

func (f *fakeGeoClient) GetTopTracks(
	_ context.Context,
	_ string,
	_ int,
	_ int,
) ([]lastfm.Track, error) {
	if f.trackErr != nil {
		return nil, f.trackErr
	}
	return f.tracks, nil
}

type fakeCarouselSearchClient struct {
	mu      sync.Mutex
	results map[string][]deezer.Track
	errs    map[string]error
}

func (f *fakeCarouselSearchClient) Find(
	_ context.Context,
	query string,
	_ uint,
) ([]deezer.Track, error) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if err := f.errs[query]; err != nil {
		return nil, err
	}
	return f.results[query], nil
}

func TestHandlerGetCarouselsMapsLastfmInvalidParameterToInvalidParam(t *testing.T) {
	apiErr, status := getCarouselErrorResponse(t, &fakeGeoClient{
		artistErr: lastfm.NewLastfmError(lastfm.ErrInvalidParameter, ""),
	})

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, domain.CodeInvalidParam, apiErr.Code)
	assert.Equal(t, "Invalid param", apiErr.Message)
	assert.Equal(t, "invalid parameter", apiErr.Cause)
}

func TestHandlerGetCarouselsMapsLastfmRateLimit(t *testing.T) {
	apiErr, status := getCarouselErrorResponse(t, &fakeGeoClient{
		artistErr: lastfm.NewLastfmError(lastfm.ErrRateLimited, ""),
	})

	assert.Equal(t, http.StatusTooManyRequests, status)
	assert.Equal(t, domain.CodeRateLimited, apiErr.Code)
	assert.Equal(t, "Rate limited", apiErr.Message)
	assert.Equal(t, "lastfm: rate limited", apiErr.Cause)
}

func TestHandlerGetCarouselsMapsLastfmProviderUnavailable(t *testing.T) {
	apiErr, status := getCarouselErrorResponse(t, &fakeGeoClient{
		artistErr: lastfm.NewLastfmError(lastfm.ErrServiceUnavailable, ""),
	})

	assert.Equal(t, http.StatusBadGateway, status)
	assert.Equal(t, domain.CodeProviderUnavailable, apiErr.Code)
	assert.Equal(t, "Provider unavailable", apiErr.Message)
	assert.Equal(t, "lastfm: service unavailable", apiErr.Cause)
}

func TestHandlerGetCarouselsMapsNoTopDataToNotFound(t *testing.T) {
	apiErr, status := getCarouselErrorResponse(t, &fakeGeoClient{})

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, domain.CodeNotFound, apiErr.Code)
	assert.Equal(t, "Not found", apiErr.Message)
	assert.Equal(t, "no top artists or tracks found", apiErr.Cause)
}

func TestHandlerGetCarouselsSkipsInvalidDeezerItemErrors(t *testing.T) {
	handler := NewHandler(NewService(
		&fakeGeoClient{
			artists: []lastfm.Artist{
				{Name: "bad artist"},
				{Name: "good artist"},
			},
		},
		&fakeCarouselSearchClient{
			results: map[string][]deezer.Track{
				"good artist": {testDeezerTrack("Good Artist", "Good Song", 101)},
			},
			errs: map[string]error{
				"bad artist": deezer.DeezerError{Kind: deezer.ErrQueryInvalid},
			},
		},
		nil,
	))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/carousels", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetCarousels(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response []CarouselResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
	require.Len(t, response, 1)
	assert.Equal(t, "Top Artists", response[0].Title)
	require.Len(t, response[0].Items, 1)
	assert.Equal(t, "Good Artist", response[0].Items[0].Name)
}

func TestHandlerGetCarouselsMapsFatalDeezerErrors(t *testing.T) {
	handler := NewHandler(NewService(
		&fakeGeoClient{
			artists: []lastfm.Artist{{Name: "artist"}},
		},
		&fakeCarouselSearchClient{
			errs: map[string]error{
				"artist": deezer.DeezerError{Kind: deezer.ErrServiceBusy},
			},
		},
		nil,
	))

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/carousels", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetCarousels(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadGateway, rec.Code)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
	assert.Equal(t, domain.CodeProviderUnavailable, apiErr.Code)
	assert.Equal(t, "Provider unavailable", apiErr.Message)
	assert.Equal(t, "deezer: service busy", apiErr.Cause)
}

func TestHandlerGetCarouselsReturnsInvalidParamForInvalidCountry(t *testing.T) {
	handler := NewHandler(NewService(&fakeGeoClient{}, &fakeCarouselSearchClient{}, nil))
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/carousels?country=mexico!", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetCarousels(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
	assert.Equal(t, domain.CodeInvalidParam, apiErr.Code)
	assert.Equal(t, "Invalid param", apiErr.Message)
	assert.Equal(t, "country must contain only letters and spaces", apiErr.Cause)
}

func TestHandlerGetCarouselsNormalizesCountryBeforeCallingService(t *testing.T) {
	geo := &fakeGeoClient{
		artists: []lastfm.Artist{{Name: "artist"}},
	}
	handler := NewHandler(NewService(
		geo,
		&fakeCarouselSearchClient{
			results: map[string][]deezer.Track{
				"artist": {testDeezerTrack("Artist", "Song", 101)},
			},
		},
		nil,
	))
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/carousels?country=%20United%20%20%20States%20", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetCarousels(c)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "united states", geo.country)
}

func getCarouselErrorResponse(t *testing.T, geo *fakeGeoClient) (domain.APIError, int) {
	t.Helper()

	handler := NewHandler(NewService(geo, &fakeCarouselSearchClient{}, nil))
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/carousels", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := handler.GetCarousels(c)
	require.NoError(t, err)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))

	return apiErr, rec.Code
}

func testDeezerTrack(artist string, title string, id int) deezer.Track {
	return deezer.Track{
		ID:    id,
		Title: title,
		Album: deezer.Album{
			ID:    id + 100,
			Title: title + " Album",
			Cover: "https://example.com/cover.jpg",
		},
		Artist: deezer.Artist{
			ID:      id + 200,
			Name:    artist,
			Picture: "https://example.com/artist.jpg",
		},
	}
}
