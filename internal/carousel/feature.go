package carousel

import (
	"github.com/labstack/echo/v5"
	"github.com/skylissh/melodihub/internal/cache"
	"github.com/skylissh/melodihub/internal/contracts"
)

type Feature struct {
	Handler *Handler
	Service *Service
}

type Options struct {
	CacheClient  *cache.Client
	GeoClient    LastfmGeoClient
	SearchClient contracts.DeezerSearchClient
}

func NewFeature(options *Options) *Feature {
	cache := NewCache(options.CacheClient)
	service := NewService(
		options.GeoClient,
		options.SearchClient,
		cache,
	)
	handler := NewHandler(service)

	return &Feature{
		Handler: handler,
		Service: service,
	}
}

func (f *Feature) RegisterRoutes(e *echo.Echo) {
	Routes(e, f.Handler)
}
