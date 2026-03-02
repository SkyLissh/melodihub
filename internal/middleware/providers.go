// Package middleware provides common middleware for the application
package middleware

import (
	"errors"

	"github.com/labstack/echo/v5"

	"github.com/skylissh/melodihub/internal/providers"
)

const providersKey = "providers"

var ErrProvidersNotFound = errors.New("providers not found")

func Providers(p *providers.Providers) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			c.Set(providersKey, p)

			return next(c)
		}
	}
}

func GetProviders(c *echo.Context) (*providers.Providers, error) {
	p, ok := c.Get(providersKey).(*providers.Providers)
	if !ok || p == nil {
		return nil, ErrProvidersNotFound
	}

	return p, nil
}
