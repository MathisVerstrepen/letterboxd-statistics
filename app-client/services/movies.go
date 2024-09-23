package services

import (
	"diikstra.fr/letterboxd-statistics/app-client/dto"
	"diikstra.fr/letterboxd-statistics/app-client/models"
)

func GetMovieInfos(movieId string, dateRange models.DateRange) (*dto.MovieInfoDTO, error) {
	movieInfoDb, err := models.Pdb.GetMovieInfos(movieId)
	if err != nil {
		return nil, err
	}

	movieViewDto, err := getViewStat(movieId, dateRange)
	if err != nil {
		return nil, err
	}

	return &dto.MovieInfoDTO{
		MovieInfoDb:  movieInfoDb,
		MovieViewDto: movieViewDto,
	}, nil
}

func getViewStat(movieId string, dateRange models.DateRange) (*dto.MovieViewDTO, error) {
	movieViewDto := dto.MovieViewDTO{}
	tsKey := models.WatchCount.TsKey(movieId, "w")

	lastMovieTs, err := models.Rdb.GetMovieLastTS(tsKey)
	if err != nil {
		return nil, err
	}

	lastValue := int(lastMovieTs.Value)
	movieViewDto.TotalViews = lastValue

	movieTsRange, err := models.Rdb.GetMovieFullRangeTS(tsKey, dateRange)
	if err != nil {
		return nil, err
	}

	startValue := int(movieTsRange[0].Value)
	movieViewDto.LastRangeViews = lastValue - startValue

	movieViewDto.Range = dateRange.ToRaw()

	return &movieViewDto, nil
}
