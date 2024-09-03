package main

import (
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"diikstra.fr/letterboxd-statistics/app-client/handlers"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))

	handlers.Init()

	e.Static("/assets", os.Getenv("ASSETS_PATH"))

	// ---- Home Routes ---- //
	e.GET("/", handlers.HomeHandler)

	// ---- Global Routes ---- //
	e.GET("/ping", handlers.GlobalPing)
	e.POST("/ping", handlers.GlobalPing)
	if os.Getenv("ENV") != "prod" {
		e.GET("/ws", handlers.GlobalHotReloadWS)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", os.Getenv("PORT"))))
}
