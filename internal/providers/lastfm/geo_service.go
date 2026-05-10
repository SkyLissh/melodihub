package lastfm

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type LastfmGeo interface {
	GetTopArtists(ctx context.Context, country string, limit int, page int) ([]Artist, error)
	GetTopTracks(ctx context.Context, country string, limit int, page int) ([]Track, error)
}

type lastfmGeo struct {
	provider *LastfmProvider
}

func (l *lastfmGeo) GetTopArtists(
	ctx context.Context,
	country string,
	limit int,
	page int,
) ([]Artist, error) {
	client := l.provider.client
	var result Response[TopArtists]

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"country": country,
			"limit":   fmt.Sprintf("%d", limit),
			"page":    fmt.Sprintf("%d", page),
		}).
		SetResult(&result).
		Get("geo.getTopArtists")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	return result.Data.TopArtists.Artists, nil
}

func (l *lastfmGeo) GetTopTracks(
	ctx context.Context,
	country string,
	limit int,
	page int,
) ([]Track, error) {
	client := l.provider.client
	var result Response[struct {
		TopTracks TopTracks `json:"tracks" validate:"required"`
	}]

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"country": country,
			"limit":   fmt.Sprintf("%d", limit),
			"page":    fmt.Sprintf("%d", page),
		}).
		SetResult(&result).
		Get("geo.getTopTracks")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	validate := validator.New()
	if err := validate.Struct(result.Data); err != nil {
		return nil, InvalidResponse(err)
	}

	return result.Data.TopTracks.Tracks, nil
}
