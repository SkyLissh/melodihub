package album

import (
	"context"

	"github.com/skylissh/melodihub/internal/cache"
)

type Cache struct {
	client *cache.Client
}

func NewCache(client *cache.Client) *Cache {
	if client == nil {
		panic("cache client cannot be nil")
	}

	return &Cache{client: client}
}

func (c *Cache) Get(ctx context.Context, id string) (*Detail, bool) {
	key := "album:" + id
	album, err := cache.GetJSON[Detail](ctx, c.client, key)
	if err != nil {
		return nil, false
	}

	return &album, true
}

func (c *Cache) Set(ctx context.Context, id string, album *Detail) bool {
	key := "album:" + id
	if err := cache.SetJSON(ctx, c.client, key, album, cache.DefaultTTL); err != nil {
		return false
	}
	return true
}
