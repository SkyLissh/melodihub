package lastfm

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type LastfmArtist interface {
	GetSimilar(ctx context.Context, name string, limit int) ([]Artist, error)
	GetDetail(ctx context.Context, name string) (*ArtistDetail, error)
	GetTopTracks(ctx context.Context, name string, limit int) (*TopTracks, error)
	GetTopAlbums(ctx context.Context, name string, limit int) (*TopAlbum, error)
}

type lastfmArtist struct {
	provider *LastfmProvider
}

func (l *lastfmArtist) GetSimilar(ctx context.Context, name string, limit int) ([]Artist, error) {
	client := l.provider.client
	var result Response[SimilarArtists]

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"artist": name,
			"limit":  fmt.Sprintf("%d", limit),
		}).
		SetResult(&result).
		Get("artist.getSimilar")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	return result.Data.SimilarArtists.Artists, nil
}

func (l *lastfmArtist) GetDetail(ctx context.Context, name string) (*ArtistDetail, error) {
	client := l.provider.client
	var result Response[struct {
		Artist ArtistDetail `json:"artist"`
	}]

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"artist": name,
		}).
		SetResult(&result).
		Get("artist.getInfo")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	return &result.Data.Artist, nil
}

func (l *lastfmArtist) GetTopTracks(ctx context.Context, name string, limit int) (*TopTracks, error) {
	client := l.provider.client
	var result Response[struct {
		TopTracks TopTracks `json:"toptracks"`
	}]

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"artist": name,
			"limit":  fmt.Sprintf("%d", limit),
		}).
		SetResult(&result).
		Get("artist.getTopTracks")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	validator := validator.New()
	if err := validator.Struct(&result.Data.TopTracks); err != nil {
		return nil, InvalidResponse(err)
	}

	return &result.Data.TopTracks, nil
}

func (l *lastfmArtist) GetTopAlbums(ctx context.Context, name string, limit int) (*TopAlbum, error) {
	client := l.provider.client
	var result Response[struct {
		TopAlbums TopAlbum `json:"topalbums"`
	}]

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"artist": name,
			"limit":  fmt.Sprintf("%d", limit),
		}).
		SetResult(&result).
		Get("artist.getTopAlbums")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	validate := validator.New()
	if err := validate.Struct(&result.Data.TopAlbums); err != nil {
		return nil, InvalidResponse(err)
	}

	return &result.Data.TopAlbums, nil
}
