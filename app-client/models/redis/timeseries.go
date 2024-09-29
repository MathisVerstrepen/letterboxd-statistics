package redis

import (
	"time"

	"diikstra.fr/letterboxd-statistics/app-client/models"
	"github.com/labstack/gommon/log"
	"github.com/redis/go-redis/v9"
)

func (rdb *DB) GetMovieFullRangeTS(ts string, fromTimestamp int64, toTimestamp int64, numberOfPoints int64) ([]redis.TSTimestampValue, error) {
	aggregationAvg := int((toTimestamp - fromTimestamp) / numberOfPoints)

	log.Infof("[REDIS] Getting %s from %d to %d", ts, fromTimestamp, toTimestamp)

	req := rdb.Client.TSRangeWithArgs(rdb.ctx, ts, int(fromTimestamp), int(toTimestamp), &redis.TSRangeOptions{
		Aggregator:     redis.Avg,
		BucketDuration: aggregationAvg,
	})
	res, err := req.Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (rdb *DB) GetMovieFullRangeTSFromNow(ts string, dateRange models.DateRange) ([]redis.TSTimestampValue, error) {
	fromTimestamp := dateRange.GetUnixTimestamp()
	toTimestamp := time.Now().UnixMilli()

	numberOfPoints := int64(1000)
	return rdb.GetMovieFullRangeTS(ts, fromTimestamp, toTimestamp, numberOfPoints)
}

func (rdb *DB) GetMovieLastTS(ts string) (redis.TSTimestampValue, error) {
	log.Infof("[REDIS] Getting latest %s", ts)

	req := rdb.Client.TSGet(rdb.ctx, ts)
	res, err := req.Result()
	if err != nil {
		return redis.TSTimestampValue{}, err
	}

	return res, nil
}
