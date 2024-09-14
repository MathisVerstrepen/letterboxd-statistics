package services

import (
	"errors"
	"os/exec"
	"time"

	"diikstra.fr/letterboxd-statistics/app-client/models"
	"github.com/labstack/gommon/log"
)

type Chart struct {
	MovieId      string
	RendererPath string
	BaseHTMLPath string
}

func (chart Chart) renderSVG(data string, metric models.Metric) ([]byte, error) {
	st := time.Now()

	lsCmd := exec.Command("phantomjs", chart.RendererPath, chart.BaseHTMLPath, data, string(metric))
	lsOut, err := lsCmd.Output()
	if err != nil {
		log.Error(err)
		return []byte{}, err
	}

	log.Infof("[CHART] Rendered in %s", time.Since(st))

	return lsOut, nil
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
		chartSvg, err := chart.renderSVG(data, metric)
		if err != nil {
			return "", err
		}

		err = models.Rdb.SetChartSVG(rdbKey, string(chartSvg))
		if err != nil {
			log.Error(err)
			return "", err
		}

		log.Info("[CHART] Chart served from rendering")
		return string(chartSvg), nil
	}

	log.Info("[CHART] Chart served from cache")
	return cachedChartSvg, nil
}
