package carousel

import (
	"context"
	"fmt"

	"github.com/skylissh/melodihub/internal/cache"
)

type Cache struct {
	client *cache.Client
}

func NewCache(client *cache.Client) *Cache {
	return &Cache{client}
}

func (c *Cache) Get(ctx context.Context, country string) ([]Section, bool) {
	key := fmt.Sprintf("carousel:%s", country)
	carousel, err := cache.GetJSON[[]Section](ctx, c.client, key)
	if err != nil {
		return nil, false
	}

	return carousel, true
}

func (c *Cache) Set(ctx context.Context, country string, carousel []Section) bool {
	key := fmt.Sprintf("carousel:%s", country)
	if err := cache.SetJSON(ctx, c.client, key, carousel, cache.DefaultTTL); err != nil {
		return false
	}
	return true
}
