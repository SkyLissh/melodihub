package lastfm

import (
	"fmt"
)

type LastfmTrack interface {
	GetTopTags(track string, artist string) (*[]Tag, error)
	GetSimilar(track string, artist string, limit int) (*[]SimilarTrack, error)
}

type lastfmTrack struct {
	provider *LastfmProvider
}

func (l *lastfmTrack) GetTopTags(track string, artist string) (*[]Tag, error) {
	client := l.provider.client

	res, err := client.R().
		SetQueryParams(map[string]string{
			"track":  track,
			"artist": artist,
		}).
		SetResult(&TopTags{}).
		Get("track.getTopTags")

	if err != nil {
		return nil, err
	}

	return &res.Result().(*TopTags).Tags, nil
}

func (l *lastfmTrack) GetSimilar(
	track string,
	artist string,
	limit int,
) (*[]SimilarTrack, error) {
	client := l.provider.client

	res, err := client.R().
		SetQueryParams(map[string]string{
			"track":  track,
			"artist": artist,
			"limit":  fmt.Sprintf("%d", limit),
		}).
		SetResult(&SimilarTracks{}).
		Get("track.getSimilar")

	if err != nil {
		return nil, err
	}

	return &res.Result().(*SimilarTracks).Tracks, nil
}
