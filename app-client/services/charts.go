package services

import (
	"errors"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	"diikstra.fr/letterboxd-statistics/app-client/models"
	"github.com/labstack/gommon/log"
)

type Chart struct {
	MovieId      string
	RendererPath string
	BaseHTMLPath string
}

func (chart Chart) renderSVG(data string, metric models.Metric, dateRange models.DateRange) ([]byte, error) {
	st := time.Now()

	lsCmd := exec.Command("phantomjs", chart.RendererPath, chart.BaseHTMLPath, data, string(metric), string(dateRange))
	lsOut, err := lsCmd.Output()
	if err != nil {
		log.Error(err)
		return []byte{}, err
	}

	log.Infof("[CHART] Rendered in %s", time.Since(st))

	return lsOut, nil
}

func optimizeSVG(chartSvg *string) {
	re := regexp.MustCompile(`-?\d+\.\d+`)

	*chartSvg = re.ReplaceAllStringFunc(*chartSvg, func(match string) string {
		f, err := strconv.ParseFloat(match, 64)
		if err != nil {
			return match
		}
		return strconv.FormatFloat(f, 'f', 1, 64)
	})
}

func (chart Chart) GetSVG(data string, metric models.Metric, dateRange models.DateRange, forceRender bool) (string, error) {
	if data == "" {
		return "", errors.New("data cannot be empty")
	}

	rdbKey := models.ChartKey(chart.MovieId, metric, dateRange)
	cachedChartSvg, err := models.Rdb.GetChartSVG(rdbKey)
	if err != nil {
		log.Error(err)
		return "", err
	}

	if cachedChartSvg == "" || forceRender {
		chartSvgBytes, err := chart.renderSVG(data, metric, dateRange)
		if err != nil {
			return "", err
		}

		chartSvg := string(chartSvgBytes)
		optimizeSVG(&chartSvg)

		err = models.Rdb.SetChartSVG(rdbKey, &chartSvg)
		if err != nil {
			log.Error(err)
			return "", err
		}

		log.Info("[CHART] Chart served from rendering")
		return chartSvg, nil
	}

	log.Info("[CHART] Chart served from cache")
	return cachedChartSvg, nil
}
