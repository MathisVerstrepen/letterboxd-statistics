package db

import (
	"context"
	"fmt"
	"log"
	"os"

	"diikstra.fr/letterboxd-statistics/app-cron/src/letterboxd"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

type RDB struct {
	Client *redis.Client
}

var Rdb RDB = RDB{}

type Metric string

const (
	WatchCount Metric = "watchcount"
	LikeCount  Metric = "likecount"
	Rating     Metric = "rating"
)

type DateRange string

const (
	Week  DateRange = "w"
	Month DateRange = "m"
	Year  DateRange = "y"
	All   DateRange = "a"
)

func (dateRange DateRange) GetUrlDateRange() string {
	switch dateRange {
	case Week:
		return "this/week"
	case Month:
		return "this/month"
	case Year:
		return "this/year"
	default:
		return ""
	}
}

func (metric Metric) TsKey(movieId string) string {
	return fmt.Sprintf("movies:%s:%s", metric, movieId)
}

func (db *RDB) Init() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Println("[INIT] WARNING : Failed to load .env file")
		log.Println(err)
	}

	db.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	res := db.Client.Conn().Ping(ctx)
	if res.Err() != nil {
		log.Default().Println("[INIT] fail to connect to redis")
		log.Fatal(res.Err().Error())
	}
}

// Check if a timeseries key name exist in redis db
// if not, create it
func (db *RDB) tsIntegrity(tsName string) error {
	ctx := context.Background()

	val, _ := db.Client.Exists(ctx, tsName).Result()
	if val != 1 {
		_, err := db.Client.TSCreate(ctx, tsName).Result()
		if err != nil {
			return err
		}
		fmt.Printf("[REDIS] Created timeseries %s\n", tsName)
	}

	return nil
}

func (db *RDB) TsAdd(tsName string, value float64) error {
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

func (db *RDB) SetPopularityOrder(moviesMeta []letterboxd.MovieMeta, dateRange DateRange) error {
	orderString := ""
	for _, movieMeta := range moviesMeta {
		orderString += movieMeta.Id + ":"
	}
	orderString = orderString[0 : len(orderString)-1]

	ctx := context.Background()
	err := db.Client.Set(ctx, fmt.Sprintf("popularityOrder:%s", dateRange), orderString, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
