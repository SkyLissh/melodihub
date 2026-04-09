package artist

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
	CacheClient           *cache.Client
	GetArtistClient       GetArtistClient
	SearchClient          contracts.DeezerSearchClient
	GetArtistDetailClient GetArtistDetailClient
}

func NewFeature(options *Options) *Feature {
	cache := NewCache(options.CacheClient)
	service := NewService(
		options.GetArtistClient,
		options.SearchClient,
		options.GetArtistDetailClient,
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
