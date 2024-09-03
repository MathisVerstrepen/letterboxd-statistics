package handlers

import (
	"net/http"

	"diikstra.fr/letterboxd-statistics/app-client/components"
	"github.com/labstack/echo/v4"
)

func HomeHandler(c echo.Context) error {
	return Render(c, http.StatusOK, components.Root(components.Home(), "Home"))
}
