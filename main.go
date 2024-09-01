package main

import (
	"fmt"

	"diikstra.fr/letterboxd-statistics/src/letterboxd"
)

func main() {
	rawFetchers := letterboxd.LoadFetcherFromProxies()

	letterboxdFetchers := letterboxd.Fetchers{}
	letterboxdFetchers.AddFetchers(rawFetchers)

	letterboxdGetter := letterboxd.LetterboxdGetter{Fetchers: letterboxdFetchers}
	popularMovies, _ := letterboxdGetter.GetPopularMovies(10, 0, letterboxd.Week)

	moviesStat, _ := letterboxdGetter.GetMovieStatsThreaded(popularMovies)
	for movieId, movieStat := range moviesStat {
		fmt.Printf("[%10s] %+v\n", movieId, *movieStat)
	}
}
