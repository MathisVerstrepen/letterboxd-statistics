package services

import (
	"fmt"
	"time"

	"diikstra.fr/letterboxd-statistics/app-client/dto"
	"diikstra.fr/letterboxd-statistics/app-client/models"
	"diikstra.fr/letterboxd-statistics/app-client/models/redis"
	"diikstra.fr/letterboxd-statistics/app-client/services/movie"
)

func GetMovieInfos(movieId string, dateRange models.DateRange) (*dto.MovieInfoDTO, error) {
	start := time.Now()

	movieInfoKey := redis.MovieInfoKey(movieId, dateRange)
	movieInfoDTO, err := redis.Rdb.GetMovieInfoDTO(movieInfoKey)
	if err == nil {
		return movieInfoDTO, nil
	}

	movieInfoDb, err := models.Pdb.GetMovieInfos(movieId)
	if err != nil {
		return nil, err
	}

	movieViewDto, err := movie.GetViewStat(movieId, dateRange)
	if err != nil {
		return nil, err
	}

	movieInfoDTO = &dto.MovieInfoDTO{
		MovieInfoDb:  movieInfoDb,
		MovieViewDto: movieViewDto,
	}
	err = redis.Rdb.SetMovieInfoDTO(movieInfoKey, movieInfoDTO)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s took %v\n", "GetMovieInfos", time.Since(start))

	return movieInfoDTO, nil
}
