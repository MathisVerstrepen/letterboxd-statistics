package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"diikstra.fr/letterboxd-statistics/app-client/components"
	"diikstra.fr/letterboxd-statistics/app-client/models"
	"diikstra.fr/letterboxd-statistics/app-client/services"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func getGraphComp(movieId string, movieMetric models.Metric) (templ.Component, error) {
	data, err := models.Rdb.GetMovieFullRangeTS(movieMetric.TsKey(movieId))
	if err != nil {
		log.Error(err)
		return nil, err
	}

	dataStrings := make([]string, len(data))
	for i, tsData := range data {
		dataStrings[i] = fmt.Sprintf("%d:%f", tsData.Timestamp, tsData.Value)
	}
	dataString := strings.Join(dataStrings, ";")

	chartEngine := services.Chart{
		MovieId:      movieId,
		RendererPath: "assets/graphs/renderHTML.js",
		BaseHTMLPath: "assets/graphs/d3.html",
	}
	svg, err := chartEngine.GetSVG(dataString, movieMetric, true)
	if err != nil {
		return nil, err
	}

	return components.Movie(svg), nil
}

func MoviePageById(c echo.Context) error {
	st := time.Now()

	movieId := c.Param("id")
	if movieId == "" {
		return errors.New("id can't be null")
	}

	movieMetricInput := c.QueryParam("metric")
	var movieMetric models.Metric
	if movieMetricInput == "listcount" {
		movieMetric = models.ListCount
	} else if movieMetricInput == "likecount" {
		movieMetric = models.LikeCount
	} else {
		movieMetric = models.WatchCount
	}

	movieInfo, err := models.Pdb.GetMovieInfos(movieId)
	if err != nil {
		return err
	}

	graphComp, err := getGraphComp(movieId, movieMetric)
	if err != nil {
		return err
	}

	return Render(c, http.StatusOK,
		components.Root(components.MovieWrapper(graphComp, movieInfo), "movie:"+movieId,
			float64(time.Since(st).Seconds()),
		),
	)
}

func GraphById(c echo.Context) error {
	movieId := c.Param("id")
	if movieId == "" {
		return errors.New("id can't be null")
	}

	movieMetricInput := c.QueryParam("metric")
	var movieMetric models.Metric
	if movieMetricInput == "listcount" {
		movieMetric = models.ListCount
	} else if movieMetricInput == "likecount" {
		movieMetric = models.LikeCount
	} else {
		movieMetric = models.WatchCount
	}

	graphComp, err := getGraphComp(movieId, movieMetric)
	if err != nil {
		return err
	}

	return Render(c, http.StatusOK, graphComp)
}
