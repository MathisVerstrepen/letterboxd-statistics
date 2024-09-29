package models

import (
	"fmt"
	"time"
)

type Metric string
type DateRange string
type LetterboxdDateRange string

const (
	WatchCount          Metric              = "watchcount"
	LikeCount           Metric              = "likecount"
	Rating              Metric              = "rating"
	LastDay             DateRange           = "d"
	LastWeek            DateRange           = "w"
	LastMonth           DateRange           = "m"
	LetterboxdLastWeek  LetterboxdDateRange = "w"
	LetterboxdLastMonth LetterboxdDateRange = "m"
	LetterboxdLastYear  LetterboxdDateRange = "y"
	LetterboxdLastAll   LetterboxdDateRange = "a"
)

func (dateRange DateRange) GetUnixTimestamp() int64 {
	timenow := time.Now()
	if dateRange == "d" {
		return timenow.Add(-1 * 24 * time.Hour).UnixMilli()
	}
	if dateRange == "w" {
		return timenow.Add(-1 * 7 * 24 * time.Hour).UnixMilli()
	}
	return timenow.Add(-1 * 30 * 7 * 24 * time.Hour).UnixMilli()
}

func (dateRange DateRange) ToString() string {
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
