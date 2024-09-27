package models

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/labstack/gommon/log"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	Client *redis.Client
	ctx    context.Context
}

var Rdb DB

type Metric string
type DateRange string
type LetterboxdDateRange string

const (
	WatchCount          Metric              = "watchcount"
	ListCount           Metric              = "listcount"
	LikeCount           Metric              = "likecount"
	LastDay             DateRange           = "d"
	LastWeek            DateRange           = "w"
	LastMonth           DateRange           = "m"
	LetterboxdLastWeek  LetterboxdDateRange = "w"
	LetterboxdLastMonth LetterboxdDateRange = "m"
	LetterboxdLastYear  LetterboxdDateRange = "y"
	LetterboxdLastAll   LetterboxdDateRange = "a"
)

func (dateRange DateRange) getUnixTimestamp() int64 {
	timenow := time.Now()
	if dateRange == "d" {
		return timenow.Add(-1 * 24 * time.Hour).UnixMilli()
	}
	if dateRange == "w" {
		return timenow.Add(-1 * 7 * 24 * time.Hour).UnixMilli()
	}
	return timenow.Add(-1 * 30 * 7 * 24 * time.Hour).UnixMilli()
}

func (dateRange DateRange) ToRaw() string {
	switch dateRange {
	case LastDay:
		return "day"
	case LastWeek:
		return "week"
	case LastMonth:
		return "month"
	}

	return ""
}

func StringToLetterboxdDateRange(rawDateRange string) LetterboxdDateRange {
	switch rawDateRange {
	case "week":
		return "w"
	case "month":
		return "m"
	case "year":
		return "y"
	case "all":
		return "a"
	default:
		return "w"
	}
}

func (metric Metric) TsKey(movieId string) string {
	return fmt.Sprintf("movies:%s:%s", metric, movieId)
}

func ChartKey(movieId string, metric Metric, dateRange DateRange) string {
	return fmt.Sprintf("movies:%s:%s:%s:chart", metric, movieId, dateRange)
}

func (db *DB) Init() {
	db.ctx = context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Warn("[INIT] WARNING : Failed to load .env file")
		log.Warn(err)
	}

	db.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	res := db.Client.Conn().Ping(db.ctx)
	if res.Err() != nil {
		log.Error("[INIT] fail to connect to redis")
		log.Fatal(res.Err().Error())
	}
}

func (rdb *DB) GetMovieFullRangeTS(ts string, dateRange DateRange) ([]redis.TSTimestampValue, error) {
	fromTimestamp := dateRange.getUnixTimestamp()
	toTimestamp := time.Now().UnixMilli()
	numberOfPoints := int64(1000)
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

func (rdb *DB) GetMovieLastTS(ts string) (redis.TSTimestampValue, error) {
	log.Infof("[REDIS] Getting latest %s", ts)

	req := rdb.Client.TSGet(rdb.ctx, ts)
	res, err := req.Result()
	if err != nil {
		return redis.TSTimestampValue{}, err
	}

	return res, nil
}

// Chart caching
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

func (rdb *DB) GetPopularityOrder(dateRange LetterboxdDateRange) ([]string, error) {
	req := rdb.Client.Get(rdb.ctx, fmt.Sprintf("popularityOrder:%s", dateRange))
	res, err := req.Result()
	if err != nil {
		return nil, err
	}

	return strings.Split(res, ":"), nil
}
