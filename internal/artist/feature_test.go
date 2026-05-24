package artist

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v5"
	"github.com/skylissh/melodihub/internal/domain"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeArtistClient struct {
	artist *deezer.Artist
	err    error
}

func (f *fakeArtistClient) GetByID(_ context.Context, _ int) (*deezer.Artist, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.artist, nil
}

type fakeArtistSearchClient struct {
	results map[string][]deezer.Track
	err     error
}

func (f *fakeArtistSearchClient) Find(_ context.Context, query string, _ uint) ([]deezer.Track, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.results[query], nil
}

type fakeArtistDetailClient struct {
	detail       *lastfm.ArtistDetail
	topTracks    *lastfm.TopTracks
	topAlbums    *lastfm.TopAlbum
	similar      []lastfm.Artist
	detailErr    error
	topTracksErr error
	topAlbumsErr error
	similarErr   error
}

func (f *fakeArtistDetailClient) GetDetail(_ context.Context, _ string) (*lastfm.ArtistDetail, error) {
	if f.detailErr != nil {
		return nil, f.detailErr
	}
	return f.detail, nil
}

func (f *fakeArtistDetailClient) GetTopTracks(_ context.Context, _ string, _ int) (*lastfm.TopTracks, error) {
	if f.topTracksErr != nil {
		return nil, f.topTracksErr
	}
	return f.topTracks, nil
}

func (f *fakeArtistDetailClient) GetTopAlbums(_ context.Context, _ string, _ int) (*lastfm.TopAlbum, error) {
	if f.topAlbumsErr != nil {
		return nil, f.topAlbumsErr
	}
	return f.topAlbums, nil
}

func (f *fakeArtistDetailClient) GetSimilar(_ context.Context, _ string, _ int) ([]lastfm.Artist, error) {
	if f.similarErr != nil {
		return nil, f.similarErr
	}
	return f.similar, nil
}

func testDeezerTrack(artistName, title string, id int) deezer.Track {
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
			Name:    artistName,
			Picture: "https://example.com/artist.jpg",
		},
	}
}

func testLastfmDetail() *lastfm.ArtistDetail {
	return &lastfm.ArtistDetail{
		Name: "Muse",
		Stats: lastfm.ArtistStats{
			Listeners: 3000000,
		},
		Bio: lastfm.ArtistBio{
			Summary: "Muse are an English rock band",
		},
	}
}

func testLastfmTopTracks() *lastfm.TopTracks {
	return &lastfm.TopTracks{
		Tracks: []lastfm.Track{
			{Name: "Hysteria", Listeners: 500000, Artist: lastfm.ArtistSimple{Name: "Muse"}},
		},
	}
}

func testLastfmTopAlbums() *lastfm.TopAlbum {
	return &lastfm.TopAlbum{
		Album: []lastfm.Album{
			{Name: "Absolution"},
		},
	}
}

func testLastfmSimilar() []lastfm.Artist {
	return []lastfm.Artist{
		{Name: "Radiohead"},
	}
}

func callGetArtist(handler *Handler, id string) (*httptest.ResponseRecorder, error) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/artist/"+id, nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPathValues(echo.PathValues{{Name: "id", Value: id}})

	err := handler.GetArtist(c)
	return rec, err
}

func getArtistErrorResponse(t *testing.T, artistClient *fakeArtistClient) (domain.APIError, int) {
	t.Helper()

	handler := NewHandler(NewService(
		artistClient,
		&fakeArtistSearchClient{},
		&fakeArtistDetailClient{},
		nil,
	))
	rec, err := callGetArtist(handler, "deezer:27")
	require.NoError(t, err)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))

	return apiErr, rec.Code
}

func TestHandlerGetArtistRejectsEmptyID(t *testing.T) {
	handler := NewHandler(NewService(
		&fakeArtistClient{},
		&fakeArtistSearchClient{},
		&fakeArtistDetailClient{},
		nil,
	))
	rec, err := callGetArtist(handler, "")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
	assert.Equal(t, domain.CodeInvalidParam, apiErr.Code)
	assert.Equal(t, "Invalid param", apiErr.Message)
	assert.Equal(t, "empty id", apiErr.Cause)
}

func TestHandlerGetArtistRejectsMissingColon(t *testing.T) {
	handler := NewHandler(NewService(
		&fakeArtistClient{},
		&fakeArtistSearchClient{},
		&fakeArtistDetailClient{},
		nil,
	))
	rec, err := callGetArtist(handler, "12345")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
	assert.Equal(t, domain.CodeInvalidParam, apiErr.Code)
	assert.Equal(t, "Invalid param", apiErr.Message)
	assert.Equal(t, "invalid id", apiErr.Cause)
}

func TestHandlerGetArtistRejectsUnknownProvider(t *testing.T) {
	handler := NewHandler(NewService(
		&fakeArtistClient{},
		&fakeArtistSearchClient{},
		&fakeArtistDetailClient{},
		nil,
	))
	rec, err := callGetArtist(handler, "spotify:123")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
	assert.Equal(t, domain.CodeInvalidParam, apiErr.Code)
	assert.Equal(t, "Invalid param", apiErr.Message)
	assert.Equal(t, "invalid provider", apiErr.Cause)
}

func TestHandlerGetArtistRejectsNonNumericID(t *testing.T) {
	handler := NewHandler(NewService(
		&fakeArtistClient{},
		&fakeArtistSearchClient{},
		&fakeArtistDetailClient{},
		nil,
	))
	rec, err := callGetArtist(handler, "deezer:abc")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
	assert.Equal(t, domain.CodeInvalidParam, apiErr.Code)
	assert.Equal(t, "Invalid param", apiErr.Message)
	assert.Equal(t, "invalid id", apiErr.Cause)
}

func TestHandlerGetArtistMapsDeezerNotFound(t *testing.T) {
	apiErr, status := getArtistErrorResponse(t, &fakeArtistClient{
		err: deezer.DeezerError{Kind: deezer.ErrNotFound},
	})

	assert.Equal(t, http.StatusNotFound, status)
	assert.Equal(t, domain.CodeNotFound, apiErr.Code)
	assert.Equal(t, "Not found", apiErr.Message)
	assert.Equal(t, "not found", apiErr.Cause)
}

func TestHandlerGetArtistMapsDeezerInvalidParam(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "parameter_missing", err: deezer.DeezerError{Kind: deezer.ErrParameterMissing}},
		{name: "invalid_parameter", err: deezer.DeezerError{Kind: deezer.ErrInvalidParameter}},
		{name: "query_invalid", err: deezer.DeezerError{Kind: deezer.ErrQueryInvalid}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr, status := getArtistErrorResponse(t, &fakeArtistClient{err: tt.err})

			assert.Equal(t, http.StatusBadRequest, status)
			assert.Equal(t, domain.CodeInvalidParam, apiErr.Code)
			assert.Equal(t, "Invalid param", apiErr.Message)
		})
	}
}

func TestHandlerGetArtistMapsDeezerRateLimited(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "quota_exceeded", err: deezer.DeezerError{Kind: deezer.ErrQuotaExceeded}},
		{name: "items_limit", err: deezer.DeezerError{Kind: deezer.ErrItemsLimit}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr, status := getArtistErrorResponse(t, &fakeArtistClient{err: tt.err})

			assert.Equal(t, http.StatusTooManyRequests, status)
			assert.Equal(t, domain.CodeRateLimited, apiErr.Code)
			assert.Equal(t, "Rate limited", apiErr.Message)
		})
	}
}

func TestHandlerGetArtistMapsDeezerProviderUnavailable(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "service_busy", err: deezer.DeezerError{Kind: deezer.ErrServiceBusy}},
		{name: "invalid_response", err: deezer.DeezerError{Kind: deezer.ErrInvalidResponse}},
		{name: "server_error", err: deezer.DeezerError{Kind: deezer.ErrServer}},
		{name: "unknown_code", err: deezer.DeezerError{Kind: deezer.ErrUnknownCode}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			apiErr, status := getArtistErrorResponse(t, &fakeArtistClient{err: tt.err})

			assert.Equal(t, http.StatusBadGateway, status)
			assert.Equal(t, domain.CodeProviderUnavailable, apiErr.Code)
			assert.Equal(t, "Provider unavailable", apiErr.Message)
		})
	}
}

func TestHandlerGetArtistMapsLastfmNotFound(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "invalid_parameter", err: lastfm.NewLastfmError(lastfm.ErrInvalidParameter, "")},
		{name: "invalid_resource", err: lastfm.NewLastfmError(lastfm.ErrInvalidResource, "")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(NewService(
				&fakeArtistClient{artist: &deezer.Artist{ID: 27, Name: "Muse", Picture: "pic"}},
				&fakeArtistSearchClient{
					results: map[string][]deezer.Track{"Muse": {testDeezerTrack("Muse", "Hysteria", 101)}},
				},
				&fakeArtistDetailClient{
					detail:    testLastfmDetail(),
					topTracks: testLastfmTopTracks(),
					topAlbums: testLastfmTopAlbums(),
					similar:   testLastfmSimilar(),
					detailErr: tt.err,
				},
				nil,
			))
			rec, err := callGetArtist(handler, "deezer:27")
			require.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, rec.Code)

			var apiErr domain.APIError
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
			assert.Equal(t, domain.CodeNotFound, apiErr.Code)
			assert.Equal(t, "Not found", apiErr.Message)
		})
	}
}

func TestHandlerGetArtistMapsLastfmRateLimited(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "rate_limited", err: lastfm.NewLastfmError(lastfm.ErrRateLimited, "")},
		{name: "too_many_requests", err: lastfm.NewLastfmError(lastfm.ErrTooManyRequests, "")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewHandler(NewService(
				&fakeArtistClient{artist: &deezer.Artist{ID: 27, Name: "Muse", Picture: "pic"}},
				&fakeArtistSearchClient{
					results: map[string][]deezer.Track{"Muse": {testDeezerTrack("Muse", "Hysteria", 101)}},
				},
				&fakeArtistDetailClient{
					detail:    testLastfmDetail(),
					topTracks: testLastfmTopTracks(),
					topAlbums: testLastfmTopAlbums(),
					similar:   testLastfmSimilar(),
					detailErr: tt.err,
				},
				nil,
			))
			rec, err := callGetArtist(handler, "deezer:27")
			require.NoError(t, err)
			assert.Equal(t, http.StatusTooManyRequests, rec.Code)

			var apiErr domain.APIError
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
			assert.Equal(t, domain.CodeRateLimited, apiErr.Code)
			assert.Equal(t, "Rate limited", apiErr.Message)
		})
	}
}

func TestHandlerGetArtistMapsLastfmProviderUnavailable(t *testing.T) {
	handler := NewHandler(NewService(
		&fakeArtistClient{artist: &deezer.Artist{ID: 27, Name: "Muse", Picture: "pic"}},
		&fakeArtistSearchClient{},
		&fakeArtistDetailClient{
			detail:    testLastfmDetail(),
			topTracks: testLastfmTopTracks(),
			topAlbums: testLastfmTopAlbums(),
			similar:   testLastfmSimilar(),
			detailErr: lastfm.NewLastfmError(lastfm.ErrServiceUnavailable, ""),
		},
		nil,
	))
	rec, err := callGetArtist(handler, "deezer:27")
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadGateway, rec.Code)

	var apiErr domain.APIError
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &apiErr))
	assert.Equal(t, domain.CodeProviderUnavailable, apiErr.Code)
	assert.Equal(t, "Provider unavailable", apiErr.Message)
}

func TestHandlerGetArtistReturnsSuccessResponse(t *testing.T) {
	handler := NewHandler(NewService(
		&fakeArtistClient{artist: &deezer.Artist{ID: 27, Name: "Muse", Picture: "https://example.com/muse.jpg"}},
		&fakeArtistSearchClient{
			results: map[string][]deezer.Track{
				"Muse": {
					testDeezerTrack("Muse", "Hysteria", 101),
				},
			},
		},
		&fakeArtistDetailClient{
			detail:    testLastfmDetail(),
			topTracks: testLastfmTopTracks(),
			topAlbums: testLastfmTopAlbums(),
			similar:   testLastfmSimilar(),
		},
		nil,
	))
	rec, err := callGetArtist(handler, "deezer:27")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var response DetailResponse
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))

	assert.Equal(t, "deezer:27", response.ID)
	assert.Equal(t, domain.ArtistType, response.Type)
	assert.Equal(t, "Muse", response.Name)
	assert.Equal(t, 3000000, response.Listeners)
	assert.Equal(t, "Muse are an English rock band", response.Bio)
}

func TestArtistErrorFromDeezerErrorMapsKnownKinds(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantCode   domain.ErrorCode
		wantStatus int
	}{
		{name: "parameter_missing", err: deezer.DeezerError{Kind: deezer.ErrParameterMissing}, wantCode: domain.CodeInvalidParam, wantStatus: http.StatusBadRequest},
		{name: "invalid_parameter", err: deezer.DeezerError{Kind: deezer.ErrInvalidParameter}, wantCode: domain.CodeInvalidParam, wantStatus: http.StatusBadRequest},
		{name: "query_invalid", err: deezer.DeezerError{Kind: deezer.ErrQueryInvalid}, wantCode: domain.CodeInvalidParam, wantStatus: http.StatusBadRequest},
		{name: "not_found", err: deezer.DeezerError{Kind: deezer.ErrNotFound}, wantCode: domain.CodeNotFound, wantStatus: http.StatusNotFound},
		{name: "quota_exceeded", err: deezer.DeezerError{Kind: deezer.ErrQuotaExceeded}, wantCode: domain.CodeRateLimited, wantStatus: http.StatusTooManyRequests},
		{name: "items_limit", err: deezer.DeezerError{Kind: deezer.ErrItemsLimit}, wantCode: domain.CodeRateLimited, wantStatus: http.StatusTooManyRequests},
		{name: "service_busy", err: deezer.DeezerError{Kind: deezer.ErrServiceBusy}, wantCode: domain.CodeProviderUnavailable, wantStatus: http.StatusBadGateway},
		{name: "permission", err: deezer.DeezerError{Kind: deezer.ErrPermission}, wantCode: domain.CodeProviderUnavailable, wantStatus: http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ArtistErrorFromDeezerError(tt.err)
			apiErr, status := domain.APIErrorFromMelodiError(result)
			assert.Equal(t, tt.wantStatus, status)
			assert.Equal(t, tt.wantCode, apiErr.Code)
		})
	}
}

func TestArtistErrorFromDeezerErrorFallsBackToInternal(t *testing.T) {
	result := ArtistErrorFromDeezerError(errors.New("unexpected transport error"))
	apiErr, status := domain.APIErrorFromMelodiError(result)

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, domain.CodeInternal, apiErr.Code)
	assert.Empty(t, apiErr.Cause)
}

func TestArtistErrorFromLastfmErrorMapsKnownKinds(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantCode   domain.ErrorCode
		wantStatus int
	}{
		{name: "invalid_parameter", err: lastfm.NewLastfmError(lastfm.ErrInvalidParameter, ""), wantCode: domain.CodeNotFound, wantStatus: http.StatusNotFound},
		{name: "invalid_resource", err: lastfm.NewLastfmError(lastfm.ErrInvalidResource, ""), wantCode: domain.CodeNotFound, wantStatus: http.StatusNotFound},
		{name: "rate_limited", err: lastfm.NewLastfmError(lastfm.ErrRateLimited, ""), wantCode: domain.CodeRateLimited, wantStatus: http.StatusTooManyRequests},
		{name: "too_many_requests", err: lastfm.NewLastfmError(lastfm.ErrTooManyRequests, ""), wantCode: domain.CodeRateLimited, wantStatus: http.StatusTooManyRequests},
		{name: "service_unavailable", err: lastfm.NewLastfmError(lastfm.ErrServiceUnavailable, ""), wantCode: domain.CodeProviderUnavailable, wantStatus: http.StatusBadGateway},
		{name: "auth_failed", err: lastfm.NewLastfmError(lastfm.ErrAuthFailed, ""), wantCode: domain.CodeProviderUnavailable, wantStatus: http.StatusBadGateway},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ArtistErrorFromLastfmError(tt.err)
			apiErr, status := domain.APIErrorFromMelodiError(result)
			assert.Equal(t, tt.wantStatus, status)
			assert.Equal(t, tt.wantCode, apiErr.Code)
		})
	}
}

func TestArtistErrorFromLastfmErrorFallsBackToInternal(t *testing.T) {
	result := ArtistErrorFromLastfmError(errors.New("unexpected transport error"))
	apiErr, status := domain.APIErrorFromMelodiError(result)

	assert.Equal(t, http.StatusInternalServerError, status)
	assert.Equal(t, domain.CodeInternal, apiErr.Code)
	assert.Empty(t, apiErr.Cause)
}