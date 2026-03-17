package main

import (
	"github.com/labstack/echo/v5"

	_ "github.com/skylissh/melodihub/docs"
	"github.com/skylissh/melodihub/internal/routes"
)

// @title			Melodi Metadata API
// @version		1.0
// @description	This is an api REST for the metadata plugin for Melodi
// @host			localhost:3000
func main() {
	e := echo.New()

	routes.CreateRoutes(e)

	if err := e.Start(":3000"); err != nil {
		e.Logger.Error("failed to start server", "error", err)
	}
}
