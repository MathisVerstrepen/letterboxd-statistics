package handlers

import (
	"net/http"
	"slices"
	"time"

	"diikstra.fr/letterboxd-statistics/app-client/components"
	"diikstra.fr/letterboxd-statistics/app-client/models"
	"github.com/labstack/echo/v4"
)

type popOrderSorter struct {
	order []string
}

func (sorter popOrderSorter) cmp(a, b models.MovieMeta) int {
	for _, movieId := range sorter.order {
		if a.Id == movieId {
			return -1
		}
		if b.Id == movieId {
			return 1
		}
	}

	return 0
}

func HomeHandler(c echo.Context) error {
	st := time.Now()

	popularityOrder, err := models.Rdb.GetPopularityOrder()
	if err != nil {
		return err
	}

	moviesMeta, err := models.Pdb.GetMoviesInfos(popularityOrder)
	if err != nil {
		return err
	}

	sorter := popOrderSorter{
		order: popularityOrder,
	}
	slices.SortFunc(moviesMeta, sorter.cmp)

	return Render(c, http.StatusOK, components.Root(components.Home(moviesMeta), "Home", float64(time.Since(st).Seconds())))
}
