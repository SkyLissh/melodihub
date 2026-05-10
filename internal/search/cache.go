package search

import (
	"context"
	"fmt"

	"github.com/skylissh/melodihub/internal/cache"
)

type Cache struct {
	client *cache.Client
}

func NewCache(client *cache.Client) *Cache {
	return &Cache{client: client}
}

func (s *Cache) Get(ctx context.Context, query Query, limit Limit) (*ResultResponse, bool) {
	key := fmt.Sprintf("search:%s:%d", query, limit)
	result, err := cache.GetJSON[ResultResponse](ctx, s.client, key)
	if err != nil {
		return nil, false
	}
	return &result, true
}

func (s *Cache) Set(
	ctx context.Context,
	query Query,
	limit Limit,
	result *ResultResponse,
) bool {
	key := fmt.Sprintf("search:%s:%d", query, limit)
	if err := cache.SetJSON(ctx, s.client, key, result, cache.DefaultTTL); err != nil {
		return false
	}
	return true
}
