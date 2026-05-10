package lastfm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArtistGetTopTracksReturnsSemanticAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "artist.getTopTracks", r.URL.Query().Get("method"))
		assert.Equal(t, "test-key", r.URL.Query().Get("api_key"))
		assert.Equal(t, "json", r.URL.Query().Get("format"))

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"error":6,"message":"Invalid parameters"}`))
	}))
	t.Cleanup(server.Close)

	provider := New("test-key")
	provider.client.SetBaseURL(server.URL)

	result, err := provider.Artist.GetTopTracks(context.Background(), "", 10)

	require.Nil(t, result)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidParameter)
}

func TestArtistGetTopTracksWrapsInvalidSuccessPayload(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"toptracks":{"track":[{"url":"https://example.test"}]}}`))
	}))
	t.Cleanup(server.Close)

	provider := New("test-key")
	provider.client.SetBaseURL(server.URL)

	result, err := provider.Artist.GetTopTracks(context.Background(), "Muse", 10)

	require.Nil(t, result)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrInvalidResponse)
}
