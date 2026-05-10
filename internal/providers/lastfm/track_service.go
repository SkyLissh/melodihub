package lastfm

import (
	"context"
	"fmt"
)

type LastfmTrack interface {
	GetTopTags(ctx context.Context, track string, artist string) ([]Tag, error)
	GetSimilar(ctx context.Context, track string, artist string, limit int) ([]SimilarTrack, error)
}

type lastfmTrack struct {
	provider *LastfmProvider
}

func (l *lastfmTrack) GetTopTags(ctx context.Context, track string, artist string) ([]Tag, error) {
	client := l.provider.client
	var result Response[struct {
		TopTags TopTags `json:"toptags"`
	}]

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"track":  track,
			"artist": artist,
		}).
		SetResult(&result).
		Get("track.getTopTags")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	return result.Data.TopTags.Tags, nil
}

func (l *lastfmTrack) GetSimilar(
	ctx context.Context,
	track string,
	artist string,
	limit int,
) ([]SimilarTrack, error) {
	client := l.provider.client
	var result Response[struct {
		SimilarTracks SimilarTracks `json:"similartracks"`
	}]

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"track":  track,
			"artist": artist,
			"limit":  fmt.Sprintf("%d", limit),
		}).
		SetResult(&result).
		Get("track.getSimilar")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	return result.Data.SimilarTracks.Tracks, nil
}
