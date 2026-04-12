package main

import (
	"github.com/labstack/echo/v5"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "github.com/skylissh/melodihub/docs"
	"github.com/skylissh/melodihub/internal/bootstrap"
	"github.com/skylissh/melodihub/internal/cache"
	"github.com/skylissh/melodihub/internal/core"
)

// @title			Melodi Metadata API
// @version		1.0
// @description	This is an api REST for the metadata plugin for Melodi
// @host			localhost:3000
func main() {
	env, err := core.NewEnv()
	if err != nil {
		panic(err)
	}

	e := echo.New()

	cacheClient, err := cache.NewClient(env.ValkeyAddr)
	if err != nil {
		panic(err)
	}

	bootstrap.Logger(e)
	bootstrap.CORS(e)

	app := bootstrap.NewApp(&bootstrap.Config{
		Env:         env,
		CacheClient: cacheClient,
	})

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/health", func(c *echo.Context) error {
		return c.JSON(200, map[string]any{"message": "ok"})
	})

	app.Search.RegisterRoutes(e)
	app.Carousel.RegisterRoutes(e)
	app.Artist.RegisterRoutes(e)

	if err := e.Start(":3000"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
