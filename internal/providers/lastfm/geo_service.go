package lastfm

import (
	"fmt"
)

type LastfmGeo interface {
	GetTopArtists(country string, limit int, page int) (*[]Artist, error)
	GetTopTracks(country string, limit int, page int) (*[]Track, error)
}

type lastfmGeo struct {
	provider *LastfmProvider
}

func (l *lastfmGeo) GetTopArtists(
	country string,
	limit int,
	page int,
) (*[]Artist, error) {
	client := l.provider.client

	res, err := client.R().
		SetQueryParams(map[string]string{
			"country": country,
			"limit":   fmt.Sprintf("%d", limit),
			"page":    fmt.Sprintf("%d", page),
		}).
		SetResult(&TopArtists{}).
		Get("geo.getTopArtists")

	if err != nil {
		return nil, err
	}

	return &res.Result().(*TopArtists).TopArtists.Artists, nil
}

func (l *lastfmGeo) GetTopTracks(
	country string,
	limit int,
	page int,
) (*[]Track, error) {
	client := l.provider.client

	res, err := client.R().
		SetQueryParams(map[string]string{
			"country": country,
			"limit":   fmt.Sprintf("%d", limit),
			"page":    fmt.Sprintf("%d", page),
		}).
		SetResult(&TopTracks{}).
		Get("geo.getTopTracks")

	if err != nil {
		return nil, err
	}

	return &res.Result().(*TopTracks).TopTracks.Tracks, nil
}
