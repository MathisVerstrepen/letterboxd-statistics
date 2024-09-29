package redis

import (
	"fmt"
	"time"

	"diikstra.fr/letterboxd-statistics/app-client/models"
)

func ChartKey(movieId string, metric models.Metric, dateRange models.DateRange) string {
	return fmt.Sprintf("movies:%s:%s:%s:chart", metric, movieId, dateRange)
}

func (rdb *DB) GetChartSVG(key string) (string, error) {
	req := rdb.Client.Get(rdb.ctx, key)
	res, err := req.Result()
	if err != nil {
		// in case of empty key, we just return an empty result but not raise the error
		if err.Error() == "redis: nil" {
			return "", nil
		}
		return "", err
	}

	return res, nil
}

func (rdb *DB) SetChartSVG(key string, svg *string) error {
	req := rdb.Client.Set(rdb.ctx, key, *svg, time.Minute*5)
	err := req.Err()

	return err
}
