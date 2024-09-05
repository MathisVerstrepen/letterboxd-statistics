package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	Client *redis.Client
	ctx    context.Context
}

var Rdb DB

type TS_NAME string

func (tsBase TS_NAME) Id(movieId string) string {
	return fmt.Sprintf("%s:%s", tsBase, movieId)
}

const (
	TS_WATCH TS_NAME = "movies:watchcount"
	TS_LIST  TS_NAME = "movies:listcount"
	TS_LIKE  TS_NAME = "movies:likecount"
)

func (db *DB) Init() {
	db.ctx = context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("[INIT] WARNING : Failed to load .env file")
		log.Println(err)
	}

	fmt.Println(os.Getenv("REDIS_HOST"))
	fmt.Println(os.Getenv("REDIS_PORT"))
	db.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: "",
		DB:       0,
	})

	res := db.Client.Conn().Ping(db.ctx)
	if res.Err() != nil {
		log.Default().Println("[INIT] fail to connect to redis")
		log.Fatal(res.Err().Error())
	}
}

func (rdb *DB) GetMovieFullRangeTS(ts string) ([]redis.TSTimestampValue, error) {
	fromTimestamp := 0
	toTimestamp := time.Now().UnixMilli()

	fmt.Printf("[REDIS] Getting %s from %d to %d\n", ts, fromTimestamp, toTimestamp)

	req := rdb.Client.TSRange(rdb.ctx, ts, fromTimestamp, int(toTimestamp))
	res, err := req.Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}
