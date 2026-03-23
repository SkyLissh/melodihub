// Package cache provides a client for interacting with the cache system.
package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/valkey-io/valkey-go"
)

var ErrCacheMiss = errors.New("cache miss")
var DefaultTTL = 1 * time.Hour

type CacheClient struct {
	client valkey.Client
}

func NewCacheClient(addr string) (*CacheClient, error) {
	client, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{addr}})
	if err != nil {
		return nil, err
	}
	return &CacheClient{client: client}, nil
}

func (c *CacheClient) Get(ctx context.Context, key string) (string, error) {
	value, err := c.client.Do(ctx, c.client.B().Get().Key(key).Build()).ToString()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return "", ErrCacheMiss
		}
		return "", fmt.Errorf("cache get error: %w", err)
	}
	return value, nil
}

func (c *CacheClient) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := c.client.Do(
		ctx, c.client.B().Set().Key(key).Value(value).Ex(ttl).Build(),
	).Error()

	if err != nil {
		return fmt.Errorf("cache set error: %w", err)
	}

	return nil
}

func (c *CacheClient) Delete(ctx context.Context, key ...string) error {
	err := c.client.Do(ctx, c.client.B().Del().Key(key...).Build()).Error()

	if err != nil {
		return fmt.Errorf("cache delete error: %w", err)
	}

	return nil
}

func GetJSON[T any](ctx context.Context, cache *CacheClient, key string) (T, error) {
	var result T

	value, err := cache.Get(ctx, key)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal([]byte(value), &result)
	if err != nil {
		return result, fmt.Errorf("cache unmarshal error: %w", err)
	}

	return result, nil
}

func SetJSON[T any](
	ctx context.Context,
	cache *CacheClient,
	key string,
	value T,
	ttl time.Duration,
) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}

	return cache.Set(ctx, key, string(jsonValue), ttl)
}
