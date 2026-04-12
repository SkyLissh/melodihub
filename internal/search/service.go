package search

import (
	"context"
	"fmt"

	"github.com/skylissh/melodihub/internal/contracts"
)

type Service struct {
	search contracts.DeezerSearchClient
	cache  *Cache
}

func NewService(search contracts.DeezerSearchClient, cache *Cache) *Service {
	return &Service{search, cache}
}

func (s *Service) GetSearch(
	ctx context.Context,
	query string,
	limit int,
) (*Result, error) {
	if s.cache != nil {
		if cached, ok := s.cache.Get(ctx, query, limit); ok {
			return cached, nil
		}
	}

	results, err := s.search.Find(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("Failed to search: %w", err)
	}

	result := Builder(results, &query)
	if s.cache != nil {
		s.cache.Set(ctx, query, limit, result)
	}

	return result, nil
}
