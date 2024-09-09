package main

import (
	"log"

	"diikstra.fr/letterboxd-statistics/app-cron/src/db"
	"diikstra.fr/letterboxd-statistics/app-cron/src/letterboxd"
)

func main() {
	db.Rdb.Init()
	db.Pdb.Init()

	rawFetchers := letterboxd.LoadFetcherFromProxies()

	letterboxdFetchers := letterboxd.Fetchers{}
	letterboxdFetchers.AddFetchers(rawFetchers)

	dbMovieIds, err := db.Pdb.GetMovieIds()
	if err != nil {
		log.Fatal(err)
	}

	letterboxdGetter := letterboxd.LetterboxdGetter{Fetchers: letterboxdFetchers}
	popularMovies, _ := letterboxdGetter.GetPopularMovies(100, 0, letterboxd.Week)

	letterboxdGetter.SetMoviePosterThreaded(popularMovies, &dbMovieIds)
	letterboxdGetter.SetMovieBackdropThreaded(popularMovies, &dbMovieIds)

	moviesStat, _ := letterboxdGetter.GetMovieStatsThreaded(popularMovies)
	for movieId, movieStat := range moviesStat {
		db.Rdb.TsAdd(db.TS_WATCH.Id(movieId), float64(movieStat.WatchCount))
		db.Rdb.TsAdd(db.TS_LIKE.Id(movieId), float64(movieStat.LikeCount))
		db.Rdb.TsAdd(db.TS_LIST.Id(movieId), float64(movieStat.ListCount))
	}

	for _, movie := range popularMovies.Movies {
		if dbMovieIds.Include(movie.Id) {
			continue
		}

		err := db.Pdb.SetMovieInfos(movie)
		if err != nil {
			log.Println(err)
		}
	}
}
