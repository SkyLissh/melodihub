package album

import (
	"context"

	"github.com/skylissh/melodihub/internal/domain"
)

type Service struct {
	cache *Cache

	album AlbumDetailClient
}

func NewService(cache *Cache, album AlbumDetailClient) *Service {
	if cache == nil {
		panic("cache cannot be nil")
	}

	return &Service{cache: cache, album: album}
}

func (s *Service) GetAlbumDetail(id *domain.ID) (*Detail, error) {
	album, err := s.album.GetDetail(context.Background(), id.RawId())
	if err != nil {
		return nil, err
	}

}
