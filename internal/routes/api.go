// Package routes provides the routes for the application.
package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"

	"github.com/rs/zerolog"

	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/skylissh/melodihub/internal/cache"
	"github.com/skylissh/melodihub/internal/core"
	"github.com/skylissh/melodihub/internal/handlers"
	"github.com/skylissh/melodihub/internal/middleware"

	"github.com/skylissh/melodihub/internal/providers"
	"github.com/skylissh/melodihub/internal/providers/deezer"
	"github.com/skylissh/melodihub/internal/providers/lastfm"
	"github.com/skylissh/melodihub/internal/providers/spotify"
)

func CreateRoutes(e *echo.Echo) {
	// Setup logger
	console := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	console.TimeLocation = time.Local

	logger := zerolog.New(console).With().Timestamp().Logger()

	e.Use(echoMiddleware.RequestLoggerWithConfig(echoMiddleware.RequestLoggerConfig{
		LogURI:     true,
		LogStatus:  true,
		LogMethod:  true,
		LogLatency: true,
		LogValuesFunc: func(c *echo.Context, v echoMiddleware.RequestLoggerValues) error {
			logger.Info().
				Str("uri", v.URI).
				Int("status", v.Status).
				Str("method", v.Method).
				Int64("latency_ms", v.Latency.Milliseconds()).
				Msg("request")

			return nil
		},
	}))

	// Setup providers
	env := core.NewEnv()

	spotifyProvider := spotify.New(env.SpotifyClientID, env.SpotifyClientSecret)
	lastfmProvider := lastfm.New(env.LastfmAPIKey)
	deezerProvider := deezer.New()

	p := &providers.Providers{
		Spotify: spotifyProvider,
		Lastfm:  lastfmProvider,
		Deezer:  deezerProvider,
	}

	e.Use(middleware.Providers(p))

	// Setup cache
	cacheClient, err := cache.NewCacheClient(env.ValkeyAddr)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to create cache client")
	}

	e.Use(middleware.Cache(cacheClient))

	// Define routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{"message": "ok"})
	})

	RegisterCarouselsRoutes(e, handlers.NewCarouselHandler())
	RegisterSearchRoutes(e, handlers.NewSearchHandler())
	RegisterArtistRoutes(e, handlers.NewArtistHandler())
}
