package dto

import "diikstra.fr/letterboxd-statistics/app-client/models"

type MovieInfoDTO struct {
	MovieInfoDb  *models.MovieMeta
	MovieViewDto *MovieViewDTO
}

type MovieViewDTO struct {
	TotalViews         int
	LastRangeViews     int
	PreviousRangeViews int
	Range              string
	Evolution          float64
}
