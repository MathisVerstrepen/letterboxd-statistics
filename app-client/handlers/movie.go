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

func getGraphComp(movieId string, movieMetric models.Metric, graphDateRange models.DateRange, letterboxdDateRange string) (templ.Component, error) {
	data, err := models.Rdb.GetMovieFullRangeTS(movieMetric.TsKey(movieId), graphDateRange)
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
	svg, err := chartEngine.GetSVG(dataString, movieMetric, graphDateRange, false)
	if err != nil {
		return nil, err
	}

	return components.Movie(svg), nil
}

func parseMetric(movieMetricInput string) models.Metric {
	var movieMetric models.Metric
	if movieMetricInput == "rating" {
		movieMetric = models.Rating
	} else if movieMetricInput == "likecount" {
		movieMetric = models.LikeCount
	} else {
		movieMetric = models.WatchCount
	}

	return movieMetric
}

func parseDateRange(dateRangeInput string) models.DateRange {
	var dateRange models.DateRange
	if dateRangeInput == "day" {
		dateRange = models.LastDay
	} else if dateRangeInput == "month" {
		dateRange = models.LastMonth
	} else {
		dateRange = models.LastWeek
	}

	return dateRange
}

func MoviePageById(c echo.Context) error {
	st := time.Now()

	movieId := c.Param("id")
	if movieId == "" {
		return errors.New("id can't be null")
	}

	movieMetric := parseMetric(c.QueryParam("metric"))
	dateRange := parseDateRange(c.QueryParam("range"))

	movieInfoDto, err := services.GetMovieInfos(movieId, dateRange)
	if err != nil {
		return err
	}

	graphComp, err := getGraphComp(movieId, movieMetric, dateRange, "w")
	if err != nil {
		return err
	}

	return Render(c, http.StatusOK,
		components.Root(components.MovieWrapper(graphComp, movieInfoDto), movieInfoDto.MovieInfoDb.Title,
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
	movieMetric := parseMetric(movieMetricInput)
	dateRangeInput := c.QueryParam("range")
	dateRange := parseDateRange(dateRangeInput)

	graphComp, err := getGraphComp(movieId, movieMetric, dateRange, "w")
	if err != nil {
		return err
	}

	movieInfoDto, err := services.GetMovieInfos(movieId, dateRange)
	if err != nil {
		return err
	}

	c.Response().Header().Set("HX-Push-Url", fmt.Sprintf("/movie/%s?metric=%s&range=%s", movieId, movieMetricInput, dateRangeInput))

	Render(c, http.StatusOK, components.StatBox(
		components.StatBoxView(movieInfoDto.MovieViewDto), "stat-1", "true"),
	)
	return Render(c, http.StatusOK, graphComp)
}
