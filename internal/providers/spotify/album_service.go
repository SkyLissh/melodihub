package spotify

import "context"

type SpotifyAlbum interface {
	GetAlbum(ctx context.Context, id string, market string) (*Album, error)
}

type spotifyAlbum struct {
	provider *SpotifyProvider
}

func (s *spotifyAlbum) GetAlbum(ctx context.Context, id string, market string) (*Album, error) {
	client := s.provider.client
	album := &Album{}

	_, err := client.R().
		SetContext(ctx).
		SetPathParam("id", id).
		SetQueryParam("market", market).
		SetResult(album).
		Get("albums/{id}")

	if err != nil {
		return nil, err
	}

	return album, nil
}
