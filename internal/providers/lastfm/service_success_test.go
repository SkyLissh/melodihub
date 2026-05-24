package lastfm

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArtistServicesReturnDecodedSuccessfulResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-key", r.URL.Query().Get("api_key"))
		assert.Equal(t, "json", r.URL.Query().Get("format"))

		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Query().Get("method") {
		case "artist.getInfo":
			assert.Equal(t, "Muse", r.URL.Query().Get("artist"))
			_, _ = w.Write([]byte(`{
				"artist": {
					"name": "Muse",
					"mbid": "",
					"url": "https://www.last.fm/music/Muse",
					"streamable": "0",
					"ontour": "0",
					"stats": {
						"listeners": "1000",
						"playcount": "2000"
					},
					"similar": {
						"artist": []
					},
					"tags": {
						"tag": []
					},
					"bio": {
						"published": "01 Jan 2020",
						"summary": "Muse summary",
						"content": "Muse content"
					}
				}
			}`))
		case "artist.getTopTracks":
			assert.Equal(t, "10", r.URL.Query().Get("limit"))
			_, _ = w.Write([]byte(`{
				"toptracks": {
					"track": [
						{
							"name": "Hysteria",
							"duration": "227",
							"mbid": "",
							"url": "https://www.last.fm/music/Muse/_/Hysteria",
							"artist": {
								"name": "Muse",
								"mbid": "",
								"url": "https://www.last.fm/music/Muse"
							},
							"listeners": "100"
						}
					]
				}
			}`))
		case "artist.getTopAlbums":
			_, _ = w.Write([]byte(`{
				"topalbums": {
					"album": [
						{
							"name": "Absolution",
							"mbid": "",
							"playcount": 100,
							"artist": {
								"name": "Muse",
								"mbid": "",
								"url": "https://www.last.fm/music/Muse"
							}
						}
					]
				}
			}`))
		case "artist.getSimilar":
			_, _ = w.Write([]byte(`{
				"similarartists": {
					"artist": [
						{
							"name": "Radiohead",
							"listeners": "500",
							"mbid": "",
							"url": "https://www.last.fm/music/Radiohead",
							"streamable": "0"
						}
					]
				}
			}`))
		default:
			t.Fatalf("unexpected method %q", r.URL.Query().Get("method"))
		}
	}))
	t.Cleanup(server.Close)

	provider := New("test-key")
	provider.client.SetBaseURL(server.URL)

	detail, err := provider.Artist.GetDetail(context.Background(), "Muse")
	require.NoError(t, err)
	assert.Equal(t, "Muse", detail.Name)
	assert.Equal(t, 1000, detail.Stats.Listeners)

	topTracks, err := provider.Artist.GetTopTracks(context.Background(), "Muse", 10)
	require.NoError(t, err)
	require.Len(t, topTracks.Tracks, 1)
	assert.Equal(t, "Hysteria", topTracks.Tracks[0].Name)

	topAlbums, err := provider.Artist.GetTopAlbums(context.Background(), "Muse", 10)
	require.NoError(t, err)
	require.Len(t, topAlbums.Album, 1)
	assert.Equal(t, "Absolution", topAlbums.Album[0].Name)

	similar, err := provider.Artist.GetSimilar(context.Background(), "Muse", 1)
	require.NoError(t, err)
	require.Len(t, similar, 1)
	assert.Equal(t, "Radiohead", similar[0].Name)
}

func TestGeoServicesReturnDecodedSuccessfulResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-key", r.URL.Query().Get("api_key"))
		assert.Equal(t, "Mexico", r.URL.Query().Get("country"))
		assert.Equal(t, "10", r.URL.Query().Get("limit"))
		assert.Equal(t, "1", r.URL.Query().Get("page"))

		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Query().Get("method") {
		case "geo.getTopArtists":
			_, _ = w.Write([]byte(`{
				"topartists": {
					"artist": [
						{
							"name": "Daft Punk",
							"listeners": "200",
							"mbid": "",
							"url": "https://www.last.fm/music/Daft+Punk",
							"streamable": "0"
						}
					]
				}
			}`))
		case "geo.getTopTracks":
			_, _ = w.Write([]byte(`{
				"tracks": {
					"track": [
						{
							"name": "One More Time",
							"duration": "320",
							"mbid": "",
							"url": "https://www.last.fm/music/Daft+Punk/_/One+More+Time",
							"artist": {
								"name": "Daft Punk",
								"mbid": "",
								"url": "https://www.last.fm/music/Daft+Punk"
							},
							"listeners": "300"
						}
					]
				}
			}`))
		default:
			t.Fatalf("unexpected method %q", r.URL.Query().Get("method"))
		}
	}))
	t.Cleanup(server.Close)

	provider := New("test-key")
	provider.client.SetBaseURL(server.URL)

	artists, err := provider.Geo.GetTopArtists(context.Background(), "Mexico", 10, 1)
	require.NoError(t, err)
	require.Len(t, artists, 1)
	assert.Equal(t, "Daft Punk", artists[0].Name)

	tracks, err := provider.Geo.GetTopTracks(context.Background(), "Mexico", 10, 1)
	require.NoError(t, err)
	require.Len(t, tracks, 1)
	assert.Equal(t, "One More Time", tracks[0].Name)
}

func TestTrackAndAlbumServicesReturnDecodedSuccessfulResponses(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-key", r.URL.Query().Get("api_key"))

		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Query().Get("method") {
		case "track.getTopTags":
			_, _ = w.Write([]byte(`{
				"toptags": {
					"tag": [
						{
							"name": "rock",
							"url": "https://www.last.fm/tag/rock",
							"count": 100
						}
					]
				}
			}`))
		case "track.getSimilar":
			_, _ = w.Write([]byte(`{
				"similartracks": {
					"track": [
						{
							"name": "Time Is Running Out",
							"duration": 236,
							"mbid": "",
							"url": "https://www.last.fm/music/Muse/_/Time+Is+Running+Out",
							"playcount": 1000,
							"artist": {
								"name": "Muse",
								"mbid": "",
								"url": "https://www.last.fm/music/Muse"
							},
							"match": 0.91
						}
					]
				}
			}`))
		case "album.getTopTags":
			_, _ = w.Write([]byte(`{
				"toptags": {
					"tag": [
						{
							"name": "alternative",
							"url": "https://www.last.fm/tag/alternative",
							"count": 50
						}
					]
				}
			}`))
		default:
			t.Fatalf("unexpected method %q", r.URL.Query().Get("method"))
		}
	}))
	t.Cleanup(server.Close)

	provider := New("test-key")
	provider.client.SetBaseURL(server.URL)

	trackTags, err := provider.Track.GetTopTags(context.Background(), "Hysteria", "Muse")
	require.NoError(t, err)
	require.Len(t, trackTags, 1)
	assert.Equal(t, "rock", trackTags[0].Name)

	similarTracks, err := provider.Track.GetSimilar(context.Background(), "Hysteria", "Muse", 1)
	require.NoError(t, err)
	require.Len(t, similarTracks, 1)
	assert.Equal(t, "Time Is Running Out", similarTracks[0].Name)

	albumTags, err := provider.Album.GetTopTags(context.Background(), "Muse", "Absolution")
	require.NoError(t, err)
	require.Len(t, albumTags, 1)
	assert.Equal(t, "alternative", albumTags[0].Name)
}
