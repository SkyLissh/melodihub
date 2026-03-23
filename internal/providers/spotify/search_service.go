package spotify

import (
	"context"
	"fmt"

	"github.com/samber/lo"

	"github.com/skylissh/melodihub/internal/utils"
)

type SpotifySearchOptions struct {
	Market string     `url:"market,omitempty"`
	Type   []InfoType `url:"type,omitempty"`
	Limit  int        `url:"limit,omitempty"`
	Offset int        `url:"offset,omitempty"`
}

type SpotifySearch interface {
	Find(ctx context.Context, query string, options *SpotifySearchOptions) (*SearchResult, error)
}

type spotifySearch struct {
	provider *SpotifyProvider
}

func (s *spotifySearch) Find(
	ctx context.Context,
	query string,
	options *SpotifySearchOptions,
) (*SearchResult, error) {
	limit := utils.Ternary(options.Limit < 0, 20, options.Limit)
	offset := utils.Ternary(options.Offset < 0, 0, options.Offset)
	market := utils.Ternary(options.Market != "", options.Market, "US")

	client := s.provider.client
	results := &SearchResult{}

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"q":      query,
			"limit":  fmt.Sprintf("%d", limit),
			"offset": fmt.Sprintf("%d", offset),
			"market": market,
		}).
		SetResult(results).
		Get("search")

	if err != nil {
		return nil, err
	}

	if results.Tracks != nil {
		results.Tracks.Items = lo.Filter(results.Tracks.Items, func(track Track, index int) bool {
			return track.Popularity > 0
		})
	}

	if results.Artists != nil {
		results.Artists.Items = lo.Filter(results.Artists.Items, func(artist Artist, index int) bool {
			return artist.Popularity > 0
		})
	}

	return results, nil
}
