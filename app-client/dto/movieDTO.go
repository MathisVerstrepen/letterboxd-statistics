package dto

import "diikstra.fr/letterboxd-statistics/app-client/models"

type MovieInfoDTO struct {
	MovieInfoDb  *models.MovieMeta
	MovieViewDto *MovieViewDTO
}

type MovieViewDTO struct {
	TotalViews     int
	LastRangeViews int
	Range          string
}
