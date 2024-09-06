package models

import (
	"context"
	"fmt"
	"os"
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

const (
	WatchCount Metric = "watchcount"
	ListCount  Metric = "listcount"
	LikeCount  Metric = "likecount"
)

func (metric Metric) TsKey(movieId string) string {
	return fmt.Sprintf("movies:%s:%s", metric, movieId)
}

func (metric Metric) ChartKey(movieId string) string {
	return fmt.Sprintf("movies:%s:%s:chart", metric, movieId)
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

func (rdb *DB) GetMovieFullRangeTS(ts string) ([]redis.TSTimestampValue, error) {
	fromTimestamp := 0
	toTimestamp := time.Now().UnixMilli()

	log.Infof("[REDIS] Getting %s from %d to %d", ts, fromTimestamp, toTimestamp)

	req := rdb.Client.TSRange(rdb.ctx, ts, fromTimestamp, int(toTimestamp))
	res, err := req.Result()
	if err != nil {
		return nil, err
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

func (rdb *DB) SetChartSVG(key string, svg string) error {
	req := rdb.Client.Set(rdb.ctx, key, svg, time.Minute*5)
	err := req.Err()

	return err
}
