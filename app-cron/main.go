package main

import (
	"fmt"

	"diikstra.fr/letterboxd-statistics/src/db"
	"diikstra.fr/letterboxd-statistics/src/letterboxd"
)

func main() {
	rdb := db.DB{}
	rdb.Init()

	rawFetchers := letterboxd.LoadFetcherFromProxies()

	letterboxdFetchers := letterboxd.Fetchers{}
	letterboxdFetchers.AddFetchers(rawFetchers)

	letterboxdGetter := letterboxd.LetterboxdGetter{Fetchers: letterboxdFetchers}
	popularMovies, _ := letterboxdGetter.GetPopularMovies(100, 0, letterboxd.Week)

	moviesStat, _ := letterboxdGetter.GetMovieStatsThreaded(popularMovies)
	for movieId, movieStat := range moviesStat {
		fmt.Printf("[%10s] %+v\n", movieId, *movieStat)

		rdb.TsAdd(db.TS_WATCH.Id(movieId), float64(movieStat.WatchCount))
		rdb.TsAdd(db.TS_LIKE.Id(movieId), float64(movieStat.LikeCount))
		rdb.TsAdd(db.TS_LIST.Id(movieId), float64(movieStat.ListCount))
	}
}
