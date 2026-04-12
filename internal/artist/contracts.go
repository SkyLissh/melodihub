package artist

import (
	"context"

	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
)

type GetArtistClient interface {
	GetByID(ctx context.Context, id int) (*deezer.Artist, error)
}

type GetArtistDetailClient interface {
	GetDetail(ctx context.Context, name string) (*lastfm.ArtistDetail, error)
	GetTopTracks(ctx context.Context, name string, limit int) (*lastfm.TopTracks, error)
	GetTopAlbums(ctx context.Context, name string, limit int) (*lastfm.TopAlbum, error)
	GetSimilar(ctx context.Context, name string, limit int) ([]lastfm.Artist, error)
}
