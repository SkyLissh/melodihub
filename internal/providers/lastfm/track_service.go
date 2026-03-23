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
	result := TopTags{}

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"track":  track,
			"artist": artist,
		}).
		SetResult(&result).
		Get("track.getTopTags")

	if err != nil {
		return nil, err
	}

	return result.Tags, nil
}

func (l *lastfmTrack) GetSimilar(
	ctx context.Context,
	track string,
	artist string,
	limit int,
) ([]SimilarTrack, error) {
	client := l.provider.client
	result := SimilarTracks{}

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
		return nil, err
	}

	return result.Tracks, nil
}
