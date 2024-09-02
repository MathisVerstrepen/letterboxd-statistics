package db

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

type DB struct {
	Client *redis.Client
}
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
	ctx := context.Background()

	db.Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	res := db.Client.Conn().Ping(ctx)
	if res.Err() != nil {
		log.Fatal(res.Err().Error())
	}
}

// Check if a timeseries key name exist in redis db
// if not, create it
func (db *DB) tsIntegrity(tsName string) error {
	ctx := context.Background()

	val, _ := db.Client.Exists(ctx, tsName).Result()
	if val != 1 {
		_, err := db.Client.TSCreate(ctx, tsName).Result()
		if err != nil {
			return err
		}
		fmt.Printf("[INIT] Created timeseries %s\n", tsName)
	}

	return nil
}

func (db *DB) TsAdd(tsName string, value float64) error {
	ctx := context.Background()

	err := db.tsIntegrity(tsName)
	if err != nil {
		return err
	}

	_, err = db.Client.TSAdd(ctx, string(tsName), "*", value).Result()
	if err != nil {
		return err
	}

	return nil
}
