package redis

import (
	"context"
	"fmt"
	"os"
	"strings"

	"diikstra.fr/letterboxd-statistics/app-client/models"
	"github.com/labstack/gommon/log"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type DB struct {
	Client *redis.Client
	ctx    context.Context
}

var Rdb DB

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

func (rdb *DB) GetPopularityOrder(dateRange models.LetterboxdDateRange) ([]string, error) {
	req := rdb.Client.Get(rdb.ctx, fmt.Sprintf("popularityOrder:%s", dateRange))
	res, err := req.Result()
	if err != nil {
		return nil, err
	}

	return strings.Split(res, ":"), nil
}
