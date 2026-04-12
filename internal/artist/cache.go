package artist

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

func (c *Cache) Get(ctx context.Context, id string) (*Detail, bool) {
	key := fmt.Sprintf("artist:%s", id)
	artist, err := cache.GetJSON[Detail](ctx, c.client, key)
	if err != nil {
		return nil, false
	}

	return &artist, true
}

func (c *Cache) Set(ctx context.Context, id string, artist *Detail) bool {
	key := fmt.Sprintf("artist:%s", id)
	if err := cache.SetJSON(ctx, c.client, key, artist, cache.DefaultTTL); err != nil {
		return false
	}
	return true
}
