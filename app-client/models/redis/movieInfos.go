package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"diikstra.fr/letterboxd-statistics/app-client/dto"
	"diikstra.fr/letterboxd-statistics/app-client/models"
)

func MovieInfoKey(movieId string, dateRange models.DateRange) string {
	return fmt.Sprintf("movies:%s:%s:info", movieId, dateRange)
}

func (rdb *DB) GetMovieInfoDTO(key string) (*dto.MovieInfoDTO, error) {
	req := rdb.Client.Get(rdb.ctx, key)
	res, err := req.Bytes()
	if err != nil {
		return nil, err
	}

	movieInfoDTO := dto.MovieInfoDTO{}
	err = json.Unmarshal(res, &movieInfoDTO)
	if err != nil {
		fmt.Print("fail to parse movieInfoDtoObject")
		return nil, err
	}

	return &movieInfoDTO, nil
}

func (rdb *DB) SetMovieInfoDTO(key string, movieInfoDTO *dto.MovieInfoDTO) error {
	data, err := json.Marshal(*movieInfoDTO)
	req := rdb.Client.Set(rdb.ctx, key, data, time.Minute*5)
	err = req.Err()

	return err
}
