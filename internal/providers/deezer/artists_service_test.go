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

func TestArtistGetByIDReturnsArtistFromSuccessfulResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/artist/27", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": 27,
			"name": "Daft Punk",
			"link": "https://www.deezer.com/artist/27",
			"picture_medium": "https://e-cdns-images.dzcdn.net/images/artist/test/250x250-000000-80-0-0.jpg"
		}`))
	}))
	t.Cleanup(server.Close)

	provider := New(validator.New())
	provider.client.SetBaseURL(server.URL)

	artist, err := provider.Artist.GetByID(context.Background(), 27)

	require.NoError(t, err)
	require.NotNil(t, artist)
	assert.Equal(t, 27, artist.ID)
	assert.Equal(t, "Daft Punk", artist.Name)
}

func TestArtistGetByIDReturnsSemanticAPIError(t *testing.T) {
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

	artist, err := provider.Artist.GetByID(context.Background(), 27)

	require.Nil(t, artist)
	require.Error(t, err)
	assert.True(t, errors.Is(err, ErrNotFound))
}
