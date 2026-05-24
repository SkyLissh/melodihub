package deezer

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skylissh/melodihub/internal/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchFindReturnsTracksFromSuccessfulResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/search", r.URL.Path)
		assert.Equal(t, "daft punk", r.URL.Query().Get("q"))
		assert.Equal(t, "1", r.URL.Query().Get("limit"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"id": 3135556,
					"title": "Harder Better Faster Stronger",
					"link": "https://www.deezer.com/track/3135556",
					"duration": 224,
					"rank": 100000,
					"explicit_lyrics": false,
					"preview": "https://cdns-preview.dzcdn.net/stream/test.mp3",
					"album": {
						"id": 302127,
						"title": "Discovery",
						"cover_medium": "https://e-cdns-images.dzcdn.net/images/cover/test/250x250-000000-80-0-0.jpg"
					},
					"artist": {
						"id": 27,
						"name": "Daft Punk",
						"link": "https://www.deezer.com/artist/27",
						"picture_medium": "https://e-cdns-images.dzcdn.net/images/artist/test/250x250-000000-80-0-0.jpg"
					}
				}
			]
		}`))
	}))
	t.Cleanup(server.Close)

	provider := New(validator.New())
	provider.client.SetBaseURL(server.URL)

	tracks, err := provider.Search.Find(context.Background(), "daft punk", 1)

	require.NoError(t, err)
	require.Len(t, tracks, 1)
	assert.Equal(t, 3135556, tracks[0].ID)
}

func TestSearchFindReturnsSemanticAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"error": {
				"type": "DataException",
				"message": "no data",
				"code": 800
			}
		}`))
	}))
	t.Cleanup(server.Close)

	provider := New(validator.New())
	provider.client.SetBaseURL(server.URL)

	tracks, err := provider.Search.Find(context.Background(), "missing", 1)

	require.Nil(t, tracks)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}

func TestSearchFindWrapsInvalidSuccessPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"data": [
				{
					"link": "https://www.deezer.com/track/3135556",
					"duration": 224,
					"album": {
						"id": 302127,
						"title": "Discovery",
						"cover_medium": "https://e-cdns-images.dzcdn.net/images/cover/test/250x250-000000-80-0-0.jpg"
					},
					"artist": {
						"id": 27,
						"name": "Daft Punk",
						"picture_medium": "https://e-cdns-images.dzcdn.net/images/artist/test/250x250-000000-80-0-0.jpg"
					}
				}
			]
		}`))
	}))
	t.Cleanup(server.Close)

	provider := New(validator.New())
	provider.client.SetBaseURL(server.URL)

	tracks, err := provider.Search.Find(context.Background(), "invalid", 1)

	require.Nil(t, tracks)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrInvalidResponse))
}
