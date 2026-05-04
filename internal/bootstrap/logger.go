package bootstrap

import (
	"os"
	"time"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func NewLogger() zerolog.Logger {
	console := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	}
	console.TimeLocation = time.Local

	return zerolog.New(console).With().Timestamp().Logger()
}

func SetGlobalLogger(logger zerolog.Logger) {
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = logger
}

func Logger(e *echo.Echo, logger zerolog.Logger) {

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

}
