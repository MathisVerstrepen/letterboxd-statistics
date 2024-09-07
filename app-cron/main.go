package main

import (
	"diikstra.fr/letterboxd-statistics/app-cron/src/db"
	"diikstra.fr/letterboxd-statistics/app-cron/src/letterboxd"
)

func main() {
	rdb := db.DB{}
	rdb.Init()

	rawFetchers := letterboxd.LoadFetcherFromProxies()

	letterboxdFetchers := letterboxd.Fetchers{}
	letterboxdFetchers.AddFetchers(rawFetchers)

	letterboxdGetter := letterboxd.LetterboxdGetter{Fetchers: letterboxdFetchers}
	popularMovies, _ := letterboxdGetter.GetPopularMovies(10, 0, letterboxd.Week)
	letterboxdGetter.SetMoviePosterThreaded(popularMovies)
	letterboxdGetter.SetMovieBackdropThreaded(popularMovies)

	moviesStat, _ := letterboxdGetter.GetMovieStatsThreaded(popularMovies)
	for movieId, movieStat := range moviesStat {
		rdb.TsAdd(db.TS_WATCH.Id(movieId), float64(movieStat.WatchCount))
		rdb.TsAdd(db.TS_LIKE.Id(movieId), float64(movieStat.LikeCount))
		rdb.TsAdd(db.TS_LIST.Id(movieId), float64(movieStat.ListCount))
	}
}
