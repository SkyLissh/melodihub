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
	result := SimilarArtists{}

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"artist": name,
			"limit":  fmt.Sprintf("%d", limit),
		}).
		SetResult(&result).
		Get("artist.getSimilar")

	if err != nil {
		return nil, err
	}

	return result.SimilarArtists.Artists, nil
}

func (l *lastfmArtist) GetDetail(ctx context.Context, name string) (*ArtistDetail, error) {
	client := l.provider.client
	result := struct {
		Artist ArtistDetail `json:"artist"`
	}{}

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"artist": name,
		}).
		SetResult(&result).
		Get("artist.getInfo")

	if err != nil {
		return nil, err
	}

	return &result.Artist, nil
}

func (l *lastfmArtist) GetTopTracks(ctx context.Context, name string, limit int) (*TopTracks, error) {
	client := l.provider.client
	result := &struct {
		TopTracks TopTracks `json:"toptracks"`
	}{}

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"artist": name,
			"limit":  fmt.Sprintf("%d", limit),
		}).
		SetResult(result).
		Get("artist.getTopTracks")

	if err != nil {
		return nil, err
	}

	validator := validator.New()
	if err := validator.Struct(&result.TopTracks); err != nil {
		return nil, err
	}

	return &result.TopTracks, nil
}

func (l *lastfmArtist) GetTopAlbums(ctx context.Context, name string, limit int) (*TopAlbum, error) {
	client := l.provider.client
	result := struct {
		TopAlbums TopAlbum `json:"topalbums"`
	}{}

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"artist": name,
			"limit":  fmt.Sprintf("%d", limit),
		}).
		SetResult(&result).
		Get("artist.getTopAlbums")

	if err != nil {
		return nil, err
	}

	validate := validator.New()
	if err := validate.Struct(&result.TopAlbums); err != nil {
		return nil, err
	}

	return &result.TopAlbums, nil
}
