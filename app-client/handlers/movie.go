package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"diikstra.fr/letterboxd-statistics/app-client/components"
	"diikstra.fr/letterboxd-statistics/app-client/models"
	"diikstra.fr/letterboxd-statistics/app-client/services"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func MoviePageById(c echo.Context) error {
	movieId := c.Param("id")
	if movieId == "" {
		return errors.New("id can't be null")
	}

	data, err := models.Rdb.GetMovieFullRangeTS(models.TS_WATCH.Id(movieId))
	if err != nil {
		log.Error(err)
		return err
	}

	dataStrings := make([]string, len(data))
	for i, tsData := range data {
		fmt.Printf("%20d - %20f\n", tsData.Timestamp, tsData.Value)
		dataStrings[i] = fmt.Sprintf("%d:%f", tsData.Timestamp, tsData.Value)
	}
	dataString := strings.Join(dataStrings, ";")

	chartEngine := services.Chart{
		MovieId:      movieId,
		RendererPath: "assets/graphs/renderHTML.js",
		BaseHTMLPath: "assets/graphs/d3.html",
	}
	svg, err := chartEngine.GetSVG(dataString)
	if err != nil {
		return err
	}

	return Render(c, http.StatusOK, components.Root(components.Movie(svg), "movie:"+movieId))
}
