package artist

import (
	"context"
	"fmt"

	"github.com/skylissh/melodihub/internal/cache"
	"github.com/skylissh/melodihub/internal/domain"
)

type Cache struct {
	client *cache.Client
}

func NewCache(client *cache.Client) *Cache {
	return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, id domain.ID) (*DetailResponse, bool) {
	key := fmt.Sprintf("artist:%s", id)
	artist, err := cache.GetJSON[DetailResponse](ctx, c.client, key)
	if err != nil {
		return nil, false
	}

	return &artist, true
}

func (c *Cache) Set(ctx context.Context, id domain.ID, artist *DetailResponse) bool {
	key := fmt.Sprintf("artist:%s", id)
	if err := cache.SetJSON(ctx, c.client, key, artist, cache.DefaultTTL); err != nil {
		return false
	}
	return true
}
