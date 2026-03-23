package lastfm

import "context"

type LastfmAlbum interface {
	GetTopTags(ctx context.Context, artistName string, albumName string) ([]Tag, error)
}

type lastfmAlbum struct {
	provider *LastfmProvider
}

func (l *lastfmAlbum) GetTopTags(
	ctx context.Context,
	artistName string,
	albumName string,
) ([]Tag, error) {
	client := l.provider.client
	result := TopTags{}

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"artist": artistName,
			"album":  albumName,
		}).
		SetResult(&result).
		Get("album.getTopTags")

	if err != nil {
		return nil, err
	}

	return result.Tags, nil
}
