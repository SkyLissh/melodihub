package bootstrap

import (
	"github.com/skylissh/melodihub/internal/artist"
	"github.com/skylissh/melodihub/internal/cache"
	"github.com/skylissh/melodihub/internal/carousel"
	"github.com/skylissh/melodihub/internal/core"
	"github.com/skylissh/melodihub/internal/search"
)

type App struct {
	Search   *search.Feature
	Carousel *carousel.Feature
	Artist   *artist.Feature
}

type Config struct {
	Env         *core.Env
	Providers   *Providers
	CacheClient *cache.Client
}

func NewApp(cfg *Config) *App {
	providers := NewProviders(cfg.Env)

	searchFeature := search.NewFeature(
		&search.Options{
			CacheClient:  cfg.CacheClient,
			SearchClient: providers.Deezer.Search,
		},
	)

	carouselFeature := carousel.NewFeature(
		&carousel.Options{
			CacheClient:  cfg.CacheClient,
			GeoClient:    cfg.Providers.Lastfm.Geo,
			SearchClient: providers.Deezer.Search,
		},
	)

	artistFeature := artist.NewFeature(
		&artist.Options{
			CacheClient:           cfg.CacheClient,
			GetArtistClient:       providers.Deezer.Artist,
			GetArtistDetailClient: providers.Lastfm.Artist,
			SearchClient:          providers.Deezer.Search,
		},
	)

	return &App{
		Search:   searchFeature,
		Carousel: carouselFeature,
		Artist:   artistFeature,
	}
}
