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
	var result Response[struct {
		TopTags TopTags `json:"toptags"`
	}]

	_, err := client.R().
		SetContext(ctx).
		SetQueryParams(map[string]string{
			"artist": artistName,
			"album":  albumName,
		}).
		SetResult(&result).
		Get("album.getTopTags")

	if err != nil {
		return nil, ServerError(err)
	}

	if result.Error != nil {
		return nil, result.Error.ToError()
	}

	return result.Data.TopTags.Tags, nil
}
