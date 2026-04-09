package carousel

import (
	"context"

	"github.com/skylissh/melodihub/internal/providers/lastfm"
)

type LastfmGeoClient interface {
	GetTopArtists(
		ctx context.Context,
		country string,
		limit int,
		page int,
	) ([]lastfm.Artist, error)

	GetTopTracks(
		ctx context.Context,
		country string,
		limit int,
		page int,
	) ([]lastfm.Track, error)
}
