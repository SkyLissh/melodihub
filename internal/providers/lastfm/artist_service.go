package lastfm

import (
	"fmt"
)

type LastfmArtist interface {
	GetSimilar(name string, limit int) (*[]Artist, error)
}

type lastfmArtist struct {
	provider *LastfmProvider
}

func (l *lastfmArtist) GetSimilar(name string, limit int) (*[]Artist, error) {
	client := l.provider.client

	res, err := client.R().
		SetQueryParams(map[string]string{
			"artist": name,
			"limit":  fmt.Sprintf("%d", limit),
		}).
		SetResult(&SimilarArtists{}).
		Get("artist.getSimilar")

	if err != nil {
		return nil, err
	}

	return &res.Result().(*SimilarArtists).SimilarArtists.Artists, nil
}
