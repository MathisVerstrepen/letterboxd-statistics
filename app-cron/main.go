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
	db.Rdb.SetPopularityOrder(popularMovies.Movies)

	letterboxdGetter.SetMoviePosterThreaded(popularMovies, &dbMovieIds)
	letterboxdGetter.SetMovieBackdropThreaded(popularMovies, &dbMovieIds)

	moviesStat, _ := letterboxdGetter.GetMovieStatsThreaded(popularMovies)
	for movieId, movieStat := range moviesStat {
		db.Rdb.TsAdd(db.WatchCount.TsKey(movieId), float64(movieStat.WatchCount))
		db.Rdb.TsAdd(db.ListCount.TsKey(movieId), float64(movieStat.ListCount))
		db.Rdb.TsAdd(db.LikeCount.TsKey(movieId), float64(movieStat.LikeCount))
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
