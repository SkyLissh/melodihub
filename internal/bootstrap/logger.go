package bootstrap

import (
	"os"
	"time"

	"github.com/labstack/echo/v5"
	echoMiddleware "github.com/labstack/echo/v5/middleware"
	"github.com/rs/zerolog"
)

func Logger(e *echo.Echo) {
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

}
