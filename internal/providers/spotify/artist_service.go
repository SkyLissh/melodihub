package spotify

import (
	"context"
	"fmt"
	"strings"

	"github.com/skylissh/melodihub/internal/utils"
)

type GetArtistAlbumsOptions struct {
	IncludeGroups []string `url:"include_groups,omitempty"`
	Limit         int      `url:"limit,omitempty"`
	Offset        int      `url:"offset,omitempty"`
	Market        string   `url:"market,omitempty"`
}

type SpotifyArtist interface {
	GetArtist(ctx context.Context, id string) (*Artist, error)
	GetManyArtist(ctx context.Context, ids []string) ([]Artist, error)
	GetAlbums(ctx context.Context, id string, options *GetArtistAlbumsOptions) (*Paginated[AlbumSimple], error)
	GetTopTracks(ctx context.Context, id string, market string) ([]Track, error)
}

type spotifyArtist struct {
	provider *SpotifyProvider
}

func (s *spotifyArtist) GetArtist(ctx context.Context, id string) (*Artist, error) {
	artist := &Artist{}
	client := s.provider.client

	_, err := client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		SetResult(artist).
		Get("artists/{id}")

	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (s *spotifyArtist) GetManyArtist(ctx context.Context, ids []string) ([]Artist, error) {
	artists := []Artist{}
	client := s.provider.client

	_, err := client.R().
		SetContext(ctx).
		SetQueryParam("ids", strings.Join(ids, ",")).
		SetResult(&artists).
		Get("artists")

	if err != nil {
		return nil, err
	}

	return artists, nil
}

func (s *spotifyArtist) GetAlbums(
	ctx context.Context,
	id string,
	options *GetArtistAlbumsOptions,
) (*Paginated[AlbumSimple], error) {
	limit := utils.Ternary(options.Limit > 0, options.Limit, 20)
	offset := utils.Ternary(options.Offset > 0, options.Offset, 0)
	market := utils.Ternary(options.Market != "", options.Market, "US")

	client := s.provider.client
	result := &Paginated[AlbumSimple]{}

	_, err := client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		SetPathParams(map[string]string{
			"limit":  fmt.Sprintf("%d", limit),
			"offset": fmt.Sprintf("%d", offset),
			"market": market,
		}).
		SetResult(result).
		Get("artists/{id}/albums")

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *spotifyArtist) GetTopTracks(ctx context.Context, id string, market string) ([]Track, error) {
	client := s.provider.client
	tracks := []Track{}

	_, err := client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		SetQueryParam("market", market).
		SetResult(&tracks).
		Get("artists/{id}/top-tracks")

	if err != nil {
		return nil, err
	}

	return tracks, nil
}
