package movie

import (
	"fmt"
	"time"

	"diikstra.fr/letterboxd-statistics/app-client/dto"
	"diikstra.fr/letterboxd-statistics/app-client/models"
	"diikstra.fr/letterboxd-statistics/app-client/models/redis"
	r "github.com/redis/go-redis/v9"
)

func GetViewStat(movieId string, dateRange models.DateRange) (*dto.MovieViewDTO, error) {
	start := time.Now()

	movieViewDto := dto.MovieViewDTO{}
	tsKey := models.WatchCount.TsKey(movieId)

	lastMovieTs, err := redis.Rdb.GetMovieLastTS(tsKey)
	if err != nil {
		return nil, err
	}

	lastPeriodEndValue := int(lastMovieTs.Value)
	movieViewDto.TotalViews = lastPeriodEndValue

	movieTsRangeFromNow, err := redis.Rdb.GetMovieFullRangeTSFromNow(tsKey, dateRange)
	if err != nil {
		return nil, err
	}

	startValue := int(movieTsRangeFromNow[0].Value)
	movieViewDto.LastRangeViews = lastPeriodEndValue - startValue

	movieTsRangePreviousPeriod, err := getMovieTsRangePreviousPeriod(dateRange, tsKey)
	if err != nil {
		return nil, err
	}

	setPreviousPeriodDTOData(&movieViewDto, &movieTsRangePreviousPeriod)

	movieViewDto.Range = dateRange.ToString()

	fmt.Printf("%s took %v\n", "GetViewStat", time.Since(start))

	return &movieViewDto, nil
}

func getMovieTsRangePreviousPeriod(dateRange models.DateRange, tsKey string) ([]r.TSTimestampValue, error) {
	endPeriod := dateRange.GetUnixTimestamp()
	startPeriod := endPeriod - (time.Now().UnixMilli() - endPeriod)

	return redis.Rdb.GetMovieFullRangeTS(tsKey, startPeriod, endPeriod, 1000)
}

func setPreviousPeriodDTOData(movieViewDto *dto.MovieViewDTO, movieTsRangePreviousPeriod *[]r.TSTimestampValue) {
	if len(*movieTsRangePreviousPeriod) == 0 {
		fmt.Println("no previous period data found !")

		movieViewDto.PreviousRangeViews = 0
		movieViewDto.Evolution = 0
	} else {
		startValue := int((*movieTsRangePreviousPeriod)[0].Value)
		lastValue := int((*movieTsRangePreviousPeriod)[len(*movieTsRangePreviousPeriod)-1].Value)
		movieViewDto.PreviousRangeViews = lastValue - startValue

		movieViewDto.Evolution = calcEvolution(movieViewDto.LastRangeViews, movieViewDto.PreviousRangeViews)
	}
}

func calcEvolution(currentValue, priorValue int) float64 {
	fmt.Printf("---- %d - %d\n", currentValue, priorValue)
	return (float64(currentValue)/float64(priorValue) - 1) * 100
}
