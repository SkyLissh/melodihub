package deezer

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/skylissh/melodihub/internal/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlbumGetDetailReturnsAlbumFromSuccessfulResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/album/302127", r.URL.Path)

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{
			"id": 302127,
			"title": "Discovery",
			"cover_medium": "https://e-cdns-images.dzcdn.net/images/cover/test/250x250-000000-80-0-0.jpg",
			"genres": {
				"data": {
					"id": 113,
					"name": "Dance",
					"picture_medium": "https://e-cdns-images.dzcdn.net/images/misc/test/250x250-000000-80-0-0.jpg"
				}
			},
			"duration": 3660,
			"fans": 1000,
			"release_date": "2001-03-12",
			"nb_tracks": 14,
			"record_type": "album",
			"explicit_lyrics": false,
			"artist": {
				"id": 27,
				"name": "Daft Punk",
				"link": "https://www.deezer.com/artist/27",
				"picture_medium": "https://e-cdns-images.dzcdn.net/images/artist/test/250x250-000000-80-0-0.jpg"
			},
			"tracks": {
				"data": {
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
			}
		}`))
	}))
	t.Cleanup(server.Close)

	provider := New(validator.New())
	provider.client.SetBaseURL(server.URL)
	service := &AlbumService{
		provider:  provider,
		validator: validator.New(),
	}

	album, err := service.GetDetail(context.Background(), 302127)

	require.NoError(t, err)
	require.NotNil(t, album)
	assert.Equal(t, 302127, album.ID)
	assert.Equal(t, "Discovery", album.Title)
	assert.Equal(t, "Daft Punk", album.Artist.Name)
	assert.Equal(t, "Harder Better Faster Stronger", album.Tracks.Data.Title)
}
