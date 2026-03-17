package lastfm

type LastfmAlbum interface {
	GetTopTags(artistName string, albumName string) ([]Tag, error)
}

type lastfmAlbum struct {
	provider *LastfmProvider
}

func (l *lastfmAlbum) GetTopTags(artistName string, albumName string) ([]Tag, error) {
	client := l.provider.client
	result := TopTags{}

	_, err := client.R().
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
