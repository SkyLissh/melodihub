package middleware

import (
	"errors"

	"github.com/labstack/echo/v5"
	"github.com/skylissh/melodihub/internal/cache"
)

const cacheKey = "cache"

var ErrCacheNotFound = errors.New("cache not found")

func Cache(client *cache.CacheClient) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set(cacheKey, client)

			return next(c)
		}
	}
}

func GetCache(c *echo.Context) (*cache.CacheClient, error) {
	p, ok := c.Get(cacheKey).(*cache.CacheClient)
	if !ok || p == nil {
		return nil, ErrCacheNotFound
	}

	return p, nil
}
