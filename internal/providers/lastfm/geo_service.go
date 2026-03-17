package lastfm

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type LastfmGeo interface {
	GetTopArtists(country string, limit int, page int) ([]Artist, error)
	GetTopTracks(country string, limit int, page int) ([]Track, error)
}

type lastfmGeo struct {
	provider *LastfmProvider
}

func (l *lastfmGeo) GetTopArtists(
	country string,
	limit int,
	page int,
) ([]Artist, error) {
	client := l.provider.client
	result := TopArtists{}

	_, err := client.R().
		SetQueryParams(map[string]string{
			"country": country,
			"limit":   fmt.Sprintf("%d", limit),
			"page":    fmt.Sprintf("%d", page),
		}).
		SetResult(&result).
		Get("geo.getTopArtists")

	if err != nil {
		return nil, err
	}

	return result.TopArtists.Artists, nil
}

func (l *lastfmGeo) GetTopTracks(
	country string,
	limit int,
	page int,
) ([]Track, error) {
	client := l.provider.client
	result := struct {
		TopTracks TopTracks `json:"tracks" validate:"required"`
	}{}

	_, err := client.R().
		SetQueryParams(map[string]string{
			"country": country,
			"limit":   fmt.Sprintf("%d", limit),
			"page":    fmt.Sprintf("%d", page),
		}).
		SetResult(&result).
		Get("geo.getTopTracks")

	if err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(result); err != nil {
		return nil, err
	}

	return result.TopTracks.Tracks, nil
}
